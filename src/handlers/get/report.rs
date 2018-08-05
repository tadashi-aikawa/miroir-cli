extern crate serde_json;

use serde_json::Value;

use clients::aws;

fn fetch_report(bucket: &String, key: &String) -> Value {
    let report_without_trials = aws::fetch_report(
        bucket,
        &format!("results/{}/report-without-trials.json", key),
    );
    let trials = aws::fetch_report(
        bucket,
        &format!("results/{}/trials.json", key),
    );
    let mut report: Value = serde_json::from_str(&report_without_trials).unwrap();
    report["trials"] = serde_json::from_str(&trials).unwrap();
    report
}

fn find_key(bucket: &String, prefix: &String) -> String {
    let keys = aws::search_keys(bucket, &format!("results/{}", prefix));

    let result = keys.into_iter()
        .filter(|x| x.contains("report-without-trials.json"))
        .collect::<Vec<String>>();
    match result.len() {
        1 => result.first().unwrap().to_string().split("/").nth(1).unwrap().to_string(),
        _n => panic!("Unable to uniquely identify key!"),
    }
}

pub fn exec(key_prefix: &String, pretty: bool) {
    let key = find_key(&"mamansoft-miroir".to_string(), key_prefix);
    let report = fetch_report(&"mamansoft-miroir".to_string(), &key);
    if pretty {
        print!("{}", serde_json::to_string_pretty(&report).unwrap());
    } else {
        print!("{}", report.to_string());
    }
}


