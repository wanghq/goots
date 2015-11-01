// Copyright 2014 The GiterLab Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// encoder for ots2
package coder

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	. "github.com/GiterLab/goots/protobuf"
)

var DebugEncoderEnable = false // 默认关闭
var DebugDecoderEnable = false // 默认关闭

func print_request_message(pb interface{}) {
	if DebugEncoderEnable {
		fmt.Println("Request Debug...")
		switch pb.(type) {
		case *CreateTableRequest:
			fmt.Println("CreateTableRequest:", proto.MarshalTextString(pb.(*CreateTableRequest)))
		case *ListTableRequest:
			fmt.Println("ListTableRequest:", proto.MarshalTextString(pb.(*ListTableRequest)))
		case *DeleteTableRequest:
			fmt.Println("DeleteTableRequest:", proto.MarshalTextString(pb.(*DeleteTableRequest)))
		case *DescribeTableRequest:
			fmt.Println("DescribeTableRequest:", proto.MarshalTextString(pb.(*DescribeTableRequest)))
		case *UpdateTableRequest:
			fmt.Println("UpdateTableRequest:", proto.MarshalTextString(pb.(*UpdateTableRequest)))
		case *GetRowRequest:
			fmt.Println("GetRowRequest:", proto.MarshalTextString(pb.(*GetRowRequest)))
		case *PutRowRequest:
			fmt.Println("PutRowRequest:", proto.MarshalTextString(pb.(*PutRowRequest)))
		case *UpdateRowRequest:
			fmt.Println("UpdateRowRequest:", proto.MarshalTextString(pb.(*UpdateRowRequest)))
		case *DeleteRowRequest:
			fmt.Println("DeleteRowRequest:", proto.MarshalTextString(pb.(*DeleteRowRequest)))
		case *BatchGetRowRequest:
			fmt.Println("BatchGetRowRequest:", proto.MarshalTextString(pb.(*BatchGetRowRequest)))
		case *BatchWriteRowRequest:
			fmt.Println("BatchWriteRowRequest:", proto.MarshalTextString(pb.(*BatchWriteRowRequest)))
		case *GetRangeRequest:
			fmt.Println("GetRangeRequest:", proto.MarshalTextString(pb.(*GetRangeRequest)))
		}
	}
}

func print_response_message(pb interface{}) {
	if DebugDecoderEnable {
		fmt.Println("Response Debug...")
		switch pb.(type) {
		case *CreateTableResponse:
			fmt.Println("CreateTableResponse:", proto.MarshalTextString(pb.(*CreateTableResponse)))
		case *DeleteTableResponse:
			fmt.Println("DeleteTableResponse:", proto.MarshalTextString(pb.(*DeleteTableResponse)))
		case *ListTableResponse:
			fmt.Println("ListTableResponse:", proto.MarshalTextString(pb.(*ListTableResponse)))
		case *UpdateTableResponse:
			fmt.Println("UpdateTableResponse:", proto.MarshalTextString(pb.(*UpdateTableResponse)))
		case *DescribeTableResponse:
			fmt.Println("DescribeTableResponse:", proto.MarshalTextString(pb.(*DescribeTableResponse)))
		case *GetRowResponse:
			fmt.Println("GetRowResponse:", proto.MarshalTextString(pb.(*GetRowResponse)))
		case *PutRowResponse:
			fmt.Println("PutRowResponse:", proto.MarshalTextString(pb.(*PutRowResponse)))
		case *UpdateRowResponse:
			fmt.Println("UpdateRowResponse:", proto.MarshalTextString(pb.(*UpdateRowResponse)))
		case *DeleteRowResponse:
			fmt.Println("DeleteRowResponse:", proto.MarshalTextString(pb.(*DeleteRowResponse)))
		case *BatchGetRowResponse:
			fmt.Println("BatchGetRowResponse:", proto.MarshalTextString(pb.(*BatchGetRowResponse)))
		case *BatchWriteRowResponse:
			fmt.Println("BatchWriteRowResponse:", proto.MarshalTextString(pb.(*BatchWriteRowResponse)))
		case *GetRangeResponse:
			fmt.Println("GetRangeResponse:", proto.MarshalTextString(pb.(*GetRangeResponse)))
		}
	}
}
