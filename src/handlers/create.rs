use clients::aws;
use RetCode;

pub fn exec(table: String, bucket: String) -> RetCode {
    eprintln!("[Start] Create bucket ---------- {:?}", bucket);
    match aws::create_bucket(&bucket) {
        Ok(_) => eprintln!("[Success] Create bucket {:?}", bucket),
        Err(e) => {
            eprintln!("[Failure] {:?}", e);
            return RetCode::FAILURE
        },
    }
    eprintln!("[End] Create bucket ---------- {:?}", bucket);

    eprintln!("[Start] Create table ---------- {:?}", table);
    match aws::create_table(&table) {
        Ok(_) => eprintln!("[Success] Create table {:?}", table),
        Err(e) => {
            eprintln!("[Failure] {:?}", e);
            return RetCode::FAILURE
        },
    }
    eprintln!("[End] Create table ---------- {:?}", table);

    RetCode::SUCCESS
}
