package database

import (
	"common/consts"
	"common/tracing"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"order/domain"
)

type MongoDatabase struct {
	db             *mongo.Client
	dbName         string
	collectionName string
}

func NewMongoDatabase(user, password, host, port, dbName, collectionName string) (db *MongoDatabase, close func(context.Context) error, err error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", user, password, host, port)
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	return &MongoDatabase{db: client, dbName: dbName, collectionName: collectionName}, client.Disconnect, nil
}

func (m *MongoDatabase) collection() *mongo.Collection {
	return m.db.Database(m.dbName).Collection(m.collectionName)
}

type orderModel struct {
	ID          bson.ObjectID      `bson:"_id"`
	OrderID     string             `bson:"order_id"`
	CustomerID  string             `bson:"customer_id"`
	Items       []*domain.Item     `bson:"items"`
	Status      consts.OrderStatus `bson:"status"`
	PaymentLink string             `bson:"payment_link"`
}

func (m *MongoDatabase) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	ctx, span := tracing.StartSpan(ctx, "Order Repository: create")
	defer span.End()

	orderID := bson.NewObjectID().Hex()
	newOrder, err := domain.NewOrder(orderID, order.CustomerID, order.Items, order.Status, order.PaymentLink)
	if err != nil {
		return nil, err
	}

	model := m.marshalOrder(newOrder)
	if _, err = m.collection().InsertOne(ctx, model); err != nil {
		return nil, err
	}

	return newOrder, nil
}

func (m *MongoDatabase) Get(ctx context.Context, orderID string, customerID string) (*domain.Order, error) {
	ctx, span := tracing.StartSpan(ctx, "Order Repository: get")
	defer span.End()

	model := &orderModel{}
	filter := bson.M{"order_id": orderID, "customer_id": customerID}
	if err := m.collection().FindOne(ctx, filter).Decode(model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &domain.NotFoundError{OrderID: orderID, CustomerID: customerID}
		}
		return nil, err
	}

	return m.unmarshalOrder(model), nil
}

func (m *MongoDatabase) Update(ctx context.Context, order *domain.Order, updateFunc func(context.Context, *domain.Order) (*domain.Order, error)) (*domain.Order, error) {
	ctx, span := tracing.StartSpan(ctx, "Order Repository: update")
	defer span.End()

	session, err := m.db.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = session.AbortTransaction(ctx)
		} else {
			_ = session.CommitTransaction(ctx)
		}
	}()

	o, err := m.Get(ctx, order.OrderID, order.CustomerID)
	if err != nil {
		return nil, err
	}

	updatedOrder, err := updateFunc(ctx, o)
	if err != nil {
		return nil, err
	}

	if _, err = m.collection().UpdateOne(ctx,
		bson.M{"order_id": updatedOrder.OrderID, "customer_id": updatedOrder.CustomerID},
		bson.M{"$set": bson.M{"status": updatedOrder.Status, "payment_link": updatedOrder.PaymentLink}},
	); err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (m *MongoDatabase) marshalOrder(order *domain.Order) *orderModel {
	id, _ := bson.ObjectIDFromHex(order.OrderID)
	return &orderModel{
		ID:          id,
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		Items:       order.Items,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
	}
}

func (m *MongoDatabase) unmarshalOrder(order *orderModel) *domain.Order {
	return &domain.Order{
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		Items:       order.Items,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
	}
}
