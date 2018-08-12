use futures::Future;
use futures::stream::Stream;
use rusoto_core::Region;
use rusoto_dynamodb::{DynamoDb, DynamoDbClient, ScanInput, DeleteItemInput, AttributeValue, DeleteItemError};
use rusoto_s3::{GetObjectRequest, ListObjectsV2Request, S3, S3Client, GetObjectError};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct Summary {
    pub hashkey: String,
    pub title: Option<String>,
    // Option is for old version
    pub begin_time: String,
}

pub fn delete_summary(table_name: String, key: String) -> Result<String, DeleteItemError> {
    let delete_key = [("hashkey".to_string(), AttributeValue { s: Some(key.to_string()), ..Default::default() })]
        .iter()
        .cloned()
        .collect::<HashMap<String, AttributeValue>>();

    let client = DynamoDbClient::simple(Region::ApNortheast1);
    let delete_item_input = DeleteItemInput {
        table_name,
        key: delete_key,
        condition_expression: Some("attribute_exists(hashkey)".to_string()),
        ..Default::default()
    };

    match client.delete_item(&delete_item_input).sync() {
        Ok(_) => Ok(format!("keys: {:?}", key)),
        Err(err) => Err(err),
    }
}

pub fn fetch_summaries(table_name: String) -> Vec<Summary> {
    let client = DynamoDbClient::simple(Region::ApNortheast1);
    let scan_input = ScanInput {
        table_name,
        ..Default::default()
    };

    match client.scan(&scan_input).sync() {
        Ok(output) => {
            let mut vec = output.items.unwrap().into_iter()
                .map(|x| {
                    Summary {
                        hashkey: x.get("hashkey").cloned().unwrap().s.unwrap(),
                        title: x.get("title").cloned().unwrap().s,
                        begin_time: x.get("begin_time").cloned().unwrap()
                            .s.unwrap()
                            .replace("/", "-"),
                    }
                })
                .collect::<Vec<Summary>>();
            vec.sort_by_key(|x| x.begin_time.clone());
            vec
        }
        Err(error) => {
            println!("Error: {:?}", error);
            panic!(error)
        }
    }
}

pub fn fetch_report(bucket: &String, key: &String) -> String {
    let client = S3Client::simple(Region::ApNortheast1);

    let get_object_request = GetObjectRequest {
        bucket: bucket.to_string(),
        key: key.to_string(),
        ..Default::default()
    };

    match client.get_object(&get_object_request).sync() {
        Ok(output) => {
            let bytes = output.body.unwrap().concat2().wait().unwrap();
            String::from_utf8(bytes).unwrap()
        }
        Err(error) => {
            println!("Error: {:?}", error);
            panic!(error)
        }
    }
}

pub fn search_keys(bucket: &String, prefix: &String) -> Vec<String> {
    let client = S3Client::simple(Region::ApNortheast1);

    let list_objects_v2_request = ListObjectsV2Request {
        bucket: bucket.to_string(),
        prefix: Some(format!("results/{}", prefix)),
        ..Default::default()
    };

    match client.list_objects_v2(&list_objects_v2_request).sync() {
        Ok(output) => {
            match output.contents {
                Some(contents) => contents.into_iter()
                    .map(|x| x.key.unwrap().clone())
                    .collect::<Vec<String>>(),
                None => vec![]
            }
        }
        Err(error) => {
            println!("Error: {:?}", error);
            panic!(error)
        }
    }
}

pub fn exists(bucket: &String, key: &String) -> Option<bool> {
    let client = S3Client::simple(Region::ApNortheast1);

    let get_object_request = GetObjectRequest {
        bucket: bucket.to_string(),
        key: format!("results/{}/report-without-trials.json", key),
        ..Default::default()
    };

    match client.get_object(&get_object_request).sync() {
        Ok(_) => Some(true),
        Err(GetObjectError::NoSuchKey(_)) => {
            Some(false)
        },
        Err(message) => {
            eprintln!("error = {:?}", message);
            None
        },
    }
}

pub fn find_key(bucket: &String, prefix: &String) -> Result<String, String> {
    let keys = search_keys(bucket, prefix);

    let result = keys.into_iter()
        .filter(|x| x.contains("report-without-trials.json"))
        .collect::<Vec<String>>();
    match result.len() {
        1 => Ok(result.first().unwrap().to_string().split("/").nth(1).unwrap().to_string()),
        _n => Err("Unable to uniquely identify key!".to_string()),
    }
}


