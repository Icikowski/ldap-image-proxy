// Package main provides the ldap-image-proxy command.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"git.sr.ht/~icikowski/ldap-image-proxy/internal/buildinfo"
	"git.sr.ht/~icikowski/ldap-image-proxy/internal/config"
	"git.sr.ht/~icikowski/ldap-image-proxy/internal/ldap"
	"git.sr.ht/~icikowski/ldap-image-proxy/internal/logs"
	"git.sr.ht/~icikowski/ldap-image-proxy/internal/web"
)

func main() {
	log := logs.InitialLogger()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	log = logs.New(cfg.Log.Level, cfg.Log.Pretty)

	ver := buildinfo.Get()
	log.Info().
		Str(logs.FieldVersion, ver.Version).
		Str(logs.FieldCommit, ver.GitReference).
		Time(logs.FieldBuildTime, ver.BuildTime).
		Str(logs.FieldGoVersion, ver.GoVersion).
		Str(logs.FieldArch, ver.Architecture).
		Msg("initializing application")

	queryBase := cfg.LDAP.QueryBase()
	ldapClient, err := ldap.NewClient(
		cfg.LDAP.URL,
		cfg.LDAP.BindUser,
		cfg.LDAP.BindPassword,
		cfg.LDAP.BaseDN,
		cfg.LDAP.AvatarAttribute,
		queryBase,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create LDAP client")
	}
	defer func() {
		if err := ldapClient.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close LDAP client")
		}
	}()

	tmpDir, err := os.MkdirTemp("", "ldap-image-proxy-*")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create temporary directory")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Error().Err(err).Msg("failed to remove temporary directory")
		}
	}()

	webHandler := web.New(
		ldapClient,
		cfg.Server.ImageTTL,
		tmpDir,
		log,
	)

	server := &http.Server{
		Addr:    cfg.Server.ServerAddr(),
		Handler: webHandler,
	}

	go func() {
		sctx := context.WithoutCancel(ctx)

		<-ctx.Done()
		log.Info().Str(logs.FieldCause, context.Cause(ctx).Error()).Msg("shutting down")

		sctx, cancel := context.WithTimeout(sctx, 5*time.Second)
		defer cancel()

		if err := server.Shutdown(sctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown server")
		}
	}()

	log.Info().Msg("started application server")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("server error")
	}
	log.Info().Msg("stopped application server")
}
