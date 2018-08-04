extern crate docopt;
extern crate futures;
extern crate rusoto_core;
extern crate rusoto_dynamodb;
extern crate rusoto_s3;
extern crate serde;
#[macro_use]
extern crate serde_derive;
extern crate serde_json;

use docopt::Docopt;
use futures::Future;
use futures::stream::Stream;
use rusoto_core::Region;
use rusoto_dynamodb::{DynamoDb, DynamoDbClient, ScanInput};
use rusoto_s3::{GetObjectRequest, S3, S3Client};
use serde_json::Value;

const USAGE: &'static str = "
Miroir CLI

Usage:
  miroir get summaries
  miroir get report <key>

Options:
  -h --help     Show this screen.
";

#[derive(Debug, Deserialize)]
struct Args {
    cmd_get: bool,
    cmd_summaries: bool,
    cmd_report: bool,
    arg_key: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Summary {
    hashkey: String,
    title: Option<String>,
    // Option is for old version
    begin_time: String,
}

fn fetch_from_s3(client: &S3Client, bucket: &String, key: &String) -> String {
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

fn fetch_report(client: &S3Client, bucket: &String, key: &String) -> String {
    let report_without_trials = fetch_from_s3(
        client,
        bucket,
        &format!("results/{}/report-without-trials.json", key),
    );
    let trials = fetch_from_s3(
        client,
        bucket,
        &format!("results/{}/trials.json", key),
    );
    let mut report: Value = serde_json::from_str(&report_without_trials).unwrap();
    report["trials"] = serde_json::from_str(&trials).unwrap();
    return report.to_string();
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

fn handle_get_report(key: &String) {
    let client = S3Client::simple(Region::ApNortheast1);
    let report = fetch_report(&client, &"mamansoft-miroir".to_string(), key);
    print!("{}", report);
}

fn handle_get_summaries() {
    let client = DynamoDbClient::simple(Region::ApNortheast1);
    let summaries = fetch_summaries(&client, "miroir".to_string());
    let output = summaries.into_iter()
        .map(|x| format!("{:30}\t{}\t{}\n", x.begin_time, &x.hashkey[0..12], x.title.unwrap()))
        .collect::<String>();
    print!("{}", output);
}

fn main() {
    let args: Args = Docopt::new(USAGE)
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());

    if args.cmd_get {
        if args.cmd_summaries {
            handle_get_summaries();
        }
        if args.cmd_report {
            handle_get_report(&args.arg_key);
        }
    }
}
