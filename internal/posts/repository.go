package posts

import (
	"context"
	"math/rand"
	"mindblog/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func col() *mongo.Collection {
	return config.DB.Collection("posts")
}

func ctx() context.Context {
	c, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return c
}

func FindPublished(page, limit int) ([]Post, int64, error) {
	filter := bson.M{"status": "published"}
	total, _ := col().CountDocuments(ctx(), filter)
	skip := int64((page - 1) * limit)
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(skip).SetLimit(int64(limit))
	cur, err := col().Find(ctx(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	var result []Post
	cur.All(ctx(), &result)
	return result, total, nil
}

func FindAll() ([]Post, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := col().Find(ctx(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var result []Post
	cur.All(ctx(), &result)
	return result, nil
}

func FindBySlug(slug string) (*Post, error) {
	var p Post
	err := col().FindOne(ctx(), bson.M{"slug": slug, "status": "published"}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func FindByID(id bson.ObjectID) (*Post, error) {
	var p Post
	err := col().FindOne(ctx(), bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func FindRandom() (*Post, error) {
	total, err := col().CountDocuments(ctx(), bson.M{"status": "published"})
	if err != nil || total == 0 {
		return nil, err
	}
	skip := rand.Int63n(total)
	opts := options.FindOne().SetSkip(skip)
	var p Post
	err = col().FindOne(ctx(), bson.M{"status": "published"}, opts).Decode(&p)
	return &p, err
}

func Create(p *Post) error {
	p.ID = bson.NewObjectID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	_, err := col().InsertOne(ctx(), p)
	return err
}

func Update(id bson.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := col().UpdateOne(ctx(), bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func Delete(id bson.ObjectID) error {
	_, err := col().DeleteOne(ctx(), bson.M{"_id": id})
	return err
}

func TogglePublish(id bson.ObjectID) (*Post, error) {
	p, err := FindByID(id)
	if err != nil {
		return nil, err
	}
	newStatus := "published"
	if p.Status == "published" {
		newStatus = "draft"
	}
	err = Update(id, bson.M{"status": newStatus})
	if err != nil {
		return nil, err
	}
	p.Status = newStatus
	return p, nil
}

func EnsureIndexes() {
	col().Indexes().CreateOne(ctx(), mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
