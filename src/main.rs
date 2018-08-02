#[macro_use]
extern crate serde_derive;

extern crate serde;
extern crate serde_json;

extern crate rusoto_core;
extern crate rusoto_dynamodb;

use rusoto_core::Region;
use rusoto_dynamodb::{DynamoDb, DynamoDbClient, ScanInput};

#[derive(Debug, Serialize, Deserialize)]
struct Summary {
    hashkey: String,
    title: Option<String>,
    // Option is for old version
    begin_time: String,
}

fn fetch_summaries(client: &DynamoDbClient, table_name: String) -> Vec<Summary> {
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

fn main() {
    let client = DynamoDbClient::simple(Region::ApNortheast1);
    let summaries = fetch_summaries(&client, "miroir".to_string());
    let output = summaries.into_iter()
        .map(|x| format!("{:30}\t{}\t{}\n", x.begin_time, &x.hashkey[0..12], x.title.unwrap()))
        .collect::<String>();
    print!("{}", output);
}
