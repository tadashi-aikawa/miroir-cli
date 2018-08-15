extern crate serde_json;

use serde_json::Value;

use clients::aws;

fn fetch_report(bucket: &String, bucket_prefix: Option<&String>, key: &String) -> Value {
    let without_trials_key = bucket_prefix.map_or(
        format!("results/{}/report-without-trials.json", key),
        |x| format!("{}/results/{}/report-without-trials.json", x, key),
    );
    let report_without_trials = aws::fetch_report(bucket, &without_trials_key);

    let trials_key = bucket_prefix.map_or(
        format!("results/{}/trials.json", key),
        |x| format!("{}/results/{}/trials.json", x, key),
    );
    let trials = aws::fetch_report(bucket, &trials_key);

    let mut report: Value = serde_json::from_str(&report_without_trials).unwrap();
    report["trials"] = serde_json::from_str(&trials).unwrap();
    report
}

pub fn exec(
    bucket_name: &String,
    bucket_prefix: Option<String>,
    key_prefix: &String,
    format: bool,
) {
    let key = aws::find_key(bucket_name, bucket_prefix.as_ref(), key_prefix);
    eprintln!("key = {:?}", key);
    match key {
        Ok(k) => {
            let report = fetch_report(bucket_name, bucket_prefix.as_ref(), &k);
            if format {
                print!("{}", serde_json::to_string_pretty(&report).unwrap());
            } else {
                print!("{}", report.to_string());
            }
        }
        Err(e) => panic!(e),
    }
}
