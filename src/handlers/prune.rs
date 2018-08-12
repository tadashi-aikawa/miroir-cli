use clients::aws;

const TABLE: &str = "miroir";
const BUCKET: &str = "mamansoft-miroir";

pub fn exec(dry: bool) {
    let pruned_keys = aws::fetch_summaries(TABLE.to_string())
        .into_iter()
        .flat_map(|x| match aws::exists(&BUCKET.to_string(), &x.hashkey.to_string()) {
            Some(true) => {
                eprintln!("[Check]: <Exists> {:?}", x.hashkey);
                None
            }
            Some(false) => {
                eprintln!("[Check]: <No report> {:?}", x.hashkey);
                Some(x.hashkey)
            }
            None => {
                eprintln!("[Check]: <Couldn't check> {:?}", x.hashkey);
                None
            }
        })
        .collect::<Vec<String>>();

    pruned_keys.into_iter()
        .for_each(|k| {
            if dry {
                eprintln!("[Dry Run]: <PRUNED> {:?}", k);
            } else {
                match aws::delete_summary(TABLE.to_string(), k) {
                    Ok(s) => eprintln!("[Run] <PRUNED>: {:?}", s),
                    Err(e) => eprintln!("[Run] <FAILURE>: {:?}", e),
                };
            }
        })
}

