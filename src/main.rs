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

mod clients;
mod handlers;

const USAGE: &'static str = "
Miroir CLI

Usage:
  miroir get summaries --table=<table>
  miroir get report <key-prefix> --bucket=<bucket> [--format]
  miroir create --table=<table> --bucket=<bucket>
  miroir prune [--dry]
  miroir --help

Options:
  -h --help            Show this screen.
  -f --format          Pretty format
  -d --dry             Dry run
  --table=<table>      DynamoDB table name
  --bucket=<bucket>    S3 bucket name
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
    flag_prefix: Option<String>,
}

fn main() {
    let args: Args = Docopt::new(USAGE)
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());

    if args.cmd_get {
        if args.cmd_summaries {
            handlers::get::summaries::exec(&args.flag_table);
        }
        if args.cmd_report {
            handlers::get::report::exec(&args.flag_bucket, &args.arg_key_prefix, args.flag_format);
        }
    } else if args.cmd_prune {
        handlers::prune::exec(args.flag_dry);
    } else if args.cmd_create {
        std::process::exit(handlers::create::exec(args.flag_table, args.flag_bucket) as i32);
    }
}
