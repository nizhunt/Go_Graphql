package database

import (
	"context"
	"go_gql/graph/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) GetJob(id string) *model.JobListing {
	jobCollection := db.client.Database("graph-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var JobListing model.JobListing
	err := jobCollection.FindOne(ctx, filter).Decode(&JobListing)
	handleErr(err)
	return &JobListing
}

func (db *DB) GetJobs() []*model.JobListing {
	jobCollection := db.client.Database("graph-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var jobListings []*model.JobListing

	// find jobs from the database
	cursor, err := jobCollection.Find(ctx, bson.D{})
	handleErr(err)

	// put the found listings stored in cursor into jobListings slice
	err = cursor.All(context.TODO(), &jobListings)
	handleErr(err)
	return jobListings
}

func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {
	jobCollection := db.client.Database("graph-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// insert the new listing into database
	inserted, err := jobCollection.InsertOne(ctx, bson.M{
		"title":       jobInfo.Title,
		"description": jobInfo.Description,
		"url":         jobInfo.URL,
		"company":     jobInfo.Company,
	})
	handleErr(err)

	// take out the newly created ID
	insertedId := inserted.InsertedID.(primitive.ObjectID).Hex()

	// return the listing back to the client along with ID
	returnJobListing := model.JobListing{
		ID:          insertedId,
		Title:       jobInfo.Title,
		Company:     jobInfo.Company,
		URL:         jobInfo.URL,
		Description: jobInfo.Description,
	}

	return &returnJobListing
}

func (db *DB) UpdateJobListing(id string, jobInfo *model.UpdateJobListingInput) *model.JobListing {
	jobCollection := db.client.Database("graph-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateJobInfo := bson.M{}

	if jobInfo.Title != nil {
		updateJobInfo["title"] = jobInfo.Title
	}
	if jobInfo.Company != nil {
		updateJobInfo["company"] = jobInfo.Company
	}
	if jobInfo.Description != nil {
		updateJobInfo["description"] = jobInfo.Description
	}
	if jobInfo.URL != nil {
		updateJobInfo["url"] = jobInfo.URL
	}

	// convert string id to object id
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}
	// find and update the data according to the above fields and return the updated bson
	results := jobCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var jobListing model.JobListing
	// convert the returned bson into a struct
	err := results.Decode(&jobListing)
	handleErr(err)
	return &jobListing
}

func (db *DB) DeleteJobListing(id string) *model.DeleteJobResponse {
	jobCollection := db.client.Database("graph-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// convert string id to object id
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	_, err := jobCollection.DeleteOne(ctx, filter)
	handleErr(err)
	return &model.DeleteJobResponse{DeleteJobID: id}
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
