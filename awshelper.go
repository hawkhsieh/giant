package main

import (
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func DynOp(input interface{}) (output interface{}, err error) {
	itype := reflect.TypeOf(input)
	Info("DynOp Name:", itype.Name())
	conf, err := InitConfig( "us-west-2" )
	Info("dynamodb_region Region", "us-west-2")
	if err != nil {
		Error("GetSchdDyn Err:", err)
		return nil, err
	}
	svc := dynamodb.New(session.New(), conf)
	switch itype.Name() {
	case "QueryInput":
		qi := input.(dynamodb.QueryInput)
		resp, operr := svc.Query(&qi)
		output = *resp
		err = operr
	case "GetItemInput":
		geti := input.(dynamodb.GetItemInput)
		resp, operr := svc.GetItem(&geti)
		output = *resp
		err = operr
	case "PutItemInput":
		puti := input.(dynamodb.PutItemInput)
		resp, operr := svc.PutItem(&puti)
		output = *resp
		err = operr
	case "DeleteItemInput":
		deli := input.(dynamodb.DeleteItemInput)
		resp, operr := svc.DeleteItem(&deli)
		output = *resp
		err = operr
	case "UpdateItemInput":
		Info("UpdateItemInput")
		updi := input.(dynamodb.UpdateItemInput)
		resp, operr := svc.UpdateItem(&updi)
		output = *resp
		err = operr
	case "BatchGetItemInput":
		Info("DynIp BatchGetItemInput")
		bti := input.(dynamodb.BatchGetItemInput)
		resp, operr := svc.BatchGetItem(&bti)
		output = *resp
		err = operr
	}
	if err != nil {
		Error("DynOp Err:", err)
		return nil, err
	}

	return output, nil

}
