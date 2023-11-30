package plugin

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type signingRoundTripper struct {
	creds      aws.Credentials
	stsClient  *sts.Client
	issuer     TokenIssuer
	assumeRole string
	signer     *v4.Signer
	base       http.RoundTripper
	region     string
	m          *sync.Mutex
}

func newSigningRoundTripper(ctx context.Context, cfg *DatasourceConfig) (*signingRoundTripper, error) {
	c, err := newConfig(ctx, config.WithRegion(cfg.STSRegion))
	if err != nil {
		return nil, err
	}
	issuer, err := newIssuer(cfg)
	if err != nil {
		return nil, err
	}
	return &signingRoundTripper{
		stsClient:  sts.NewFromConfig(c),
		issuer:     issuer,
		assumeRole: cfg.AssumeRole,
		signer:     v4.NewSigner(),
		base:       http.DefaultTransport,
		region:     cfg.MonitoringRegion,
		m:          &sync.Mutex{},
	}, nil
}

func (rt *signingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	err := rt.sign(req, rt.region)
	if err != nil {
		return nil, err
	}
	resp, err := rt.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rt *signingRoundTripper) retrieve(ctx context.Context) (aws.Credentials, error) {
	rt.m.Lock()
	defer rt.m.Unlock()

	if rt.creds.HasKeys() && !rt.creds.Expired() {
		return rt.creds, nil
	}

	token, err := rt.issuer.IssueAccessToken(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	identity, err := rt.stsClient.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(rt.assumeRole),
		WebIdentityToken: aws.String(token),
		DurationSeconds:  aws.Int32(3600),
		RoleSessionName:  aws.String(PluginID),
	})
	if err != nil {
		return aws.Credentials{}, err
	}

	rt.creds = aws.Credentials{
		AccessKeyID:     *identity.Credentials.AccessKeyId,
		SecretAccessKey: *identity.Credentials.SecretAccessKey,
		SessionToken:    *identity.Credentials.SessionToken,
		CanExpire:       true,
		Expires:         *identity.Credentials.Expiration,
	}

	return rt.creds, nil
}

func (rt *signingRoundTripper) sign(r *http.Request, signingRegion string) error {
	if r.Body == nil {
		r.Body = http.NoBody
	}
	ctx := r.Context()
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r.Body, buf)
	_, err := buf.ReadFrom(tee)
	if err != nil {
		return err
	}

	b := sha256.Sum256(buf.Bytes())
	payloadHash := hex.EncodeToString(b[:])

	r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	creds, err := rt.retrieve(ctx)
	if err != nil {
		return err
	}

	err = rt.signer.SignHTTP(ctx, creds, r, payloadHash, ServiceID, signingRegion, time.Now())
	return err
}
