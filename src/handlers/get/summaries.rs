use clients::aws;

const TABLE: &str = "miroir";

pub fn exec() {
    let summaries = aws::fetch_summaries(TABLE.to_string());
    let output = summaries.into_iter()
        .map(|x| format!("{:30}\t{}\t{}\n", x.begin_time, &x.hashkey[0..12], x.title.unwrap()))
        .collect::<String>();
    print!("{}", output);
}

