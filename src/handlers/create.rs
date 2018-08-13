use clients::aws;


pub fn exec(table: String, bucket: String, prefix: Option<String>) {
    eprintln!("table = {:?}", table);
    eprintln!("bucket = {:?}", bucket);
    eprintln!("prefix = {:?}", prefix);
    aws::create_bucket(&bucket);
}

