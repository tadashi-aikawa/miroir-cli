extern crate docopt;
extern crate futures;
extern crate regex;
extern crate rusoto_core;
extern crate rusoto_dynamodb;
extern crate rusoto_s3;
extern crate serde;
#[macro_use]
extern crate serde_derive;
extern crate openssl_probe;
extern crate serde_json;

use docopt::Docopt;

mod clients;
mod handlers;

const USAGE: &'static str = "
Miroir CLI

Usage:
  miroir get summaries --table=<table>
  miroir get report <key-prefix> --bucket=<bucket> [--bucket-prefix=<bucket-prefix>] [--format]
  miroir create --table=<table> --bucket=<bucket>
  miroir prune  --table=<table> --bucket=<bucket> [--bucket-prefix=<bucket-prefix>] [--dry]
  miroir --help

Options:
  -h --help                          Show this screen.
  -f --format                        Pretty format
  -d --dry                           Dry run
  --table=<table>                    DynamoDB table name
  --bucket=<bucket>                  S3 bucket name
  --bucket-prefix=<bucket-prefix>    S3 bucket prefix (directory)
";

pub enum RetCode {
    SUCCESS = 0,
    FAILURE = 1,
}

#[derive(Debug, Deserialize)]
struct Args {
    cmd_get: bool,
    cmd_create: bool,
    cmd_prune: bool,
    cmd_summaries: bool,
    cmd_report: bool,
    arg_key_prefix: String,
    flag_format: bool,
    flag_dry: bool,
    flag_table: String,
    flag_bucket: String,
    flag_bucket_prefix: Option<String>,
}

fn main() {
    openssl_probe::init_ssl_cert_env_vars();

    let args: Args = Docopt::new(USAGE)
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());

    if args.cmd_get {
        if args.cmd_summaries {
            handlers::get::summaries::exec(&args.flag_table);
        }
        if args.cmd_report {
            handlers::get::report::exec(
                &args.flag_bucket,
                args.flag_bucket_prefix,
                &args.arg_key_prefix,
                args.flag_format,
            );
        }
    } else if args.cmd_prune {
        handlers::prune::exec(&args.flag_table, &args.flag_bucket, args.flag_bucket_prefix, args.flag_dry);
    } else if args.cmd_create {
        std::process::exit(handlers::create::exec(args.flag_table, args.flag_bucket) as i32);
    }
}
