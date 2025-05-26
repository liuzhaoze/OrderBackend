package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"io"
	"net/http"
)

type HttpHandler struct {
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) HandleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf("error reading request body: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	stripeEndpointSecret := viper.GetString("STRIPE_ENDPOINT_SECRET")
	if stripeEndpointSecret == "" {
		logrus.Errorln("empty stripe endpoint secret, please set STRIPE_ENDPOINT_SECRET environment variable")
		return
	}

	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), stripeEndpointSecret)
	if err != nil {
		logrus.Errorf("error verifying webhook signature: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err = json.Unmarshal(event.Data.Raw, &session); err != nil {
			logrus.Errorf("error unmarshalling checkout session: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Infoln("checkout session is paid")
		}
	default:
		logrus.Warnf("unexpected event type: %v\n", event.Type)
	}

	c.JSON(http.StatusOK, nil)
}
