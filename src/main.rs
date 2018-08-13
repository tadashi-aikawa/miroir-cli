extern crate futures;
extern crate rusoto_core;
extern crate rusoto_dynamodb;
extern crate rusoto_s3;
extern crate serde;
#[macro_use]
extern crate serde_derive;
extern crate serde_json;

extern crate docopt;


use docopt::Docopt;

mod handlers;
mod clients;

const USAGE: &'static str = "
Miroir CLI

Usage:
  miroir get summaries
  miroir get report <key-prefix> [--format]
  miroir prune [--dry]
  miroir --help

Options:
  -h --help     Show this screen.
  -f --format   Pretty format
  -d --dry      Dry run
";

#[derive(Debug, Deserialize)]
struct Args {
    cmd_get: bool,
    cmd_prune: bool,
    cmd_summaries: bool,
    cmd_report: bool,
    arg_key_prefix: String,
    flag_format: bool,
    flag_dry: bool,
}

fn main() {
    let args: Args = Docopt::new(USAGE)
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());

    if args.cmd_get {
        if args.cmd_summaries {
            handlers::get::summaries::exec();
        }
        if args.cmd_report {
            handlers::get::report::exec(&args.arg_key_prefix, args.flag_format);
        }
    }

    if args.cmd_prune {
        handlers::prune::exec(args.flag_dry);
    }

}
