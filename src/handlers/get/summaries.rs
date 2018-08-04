use clients::aws;

pub fn exec() {
    let summaries = aws::fetch_summaries("miroir".to_string());
    let output = summaries.into_iter()
        .map(|x| format!("{:30}\t{}\t{}\n", x.begin_time, &x.hashkey[0..12], x.title.unwrap()))
        .collect::<String>();
    print!("{}", output);
}

