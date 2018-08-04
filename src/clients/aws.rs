use rusoto_core::Region;
use rusoto_dynamodb::{DynamoDb, DynamoDbClient, ScanInput};
use rusoto_s3::{GetObjectRequest, S3, S3Client};
use futures::Future;
use futures::stream::Stream;

#[derive(Debug, Serialize, Deserialize)]
pub struct Summary {
    pub hashkey: String,
    pub title: Option<String>,
    // Option is for old version
    pub begin_time: String,
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
