extern crate serde_json;

use serde_json::Value;

use clients::aws;

fn fetch_report(bucket: &String, key: &String) -> Value {
    let report_without_trials = aws::fetch_report(
        bucket,
        &format!("results/{}/report-without-trials.json", key),
    );
    let trials = aws::fetch_report(bucket, &format!("results/{}/trials.json", key));
    let mut report: Value = serde_json::from_str(&report_without_trials).unwrap();
    report["trials"] = serde_json::from_str(&trials).unwrap();
    report
}

pub fn exec(bucket_name: &String, key_prefix: &String, format: bool) {
    let key = aws::find_key(bucket_name, key_prefix);
    match key {
        Ok(k) => {
            let report = fetch_report(bucket_name, &k);
            if format {
                print!("{}", serde_json::to_string_pretty(&report).unwrap());
            } else {
                print!("{}", report.to_string());
            }
        }
        Err(e) => panic!(e),
    }
}
