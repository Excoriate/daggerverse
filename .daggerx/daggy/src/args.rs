use clap::Parser;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct Args {
    /// Task is the name of the task to run
    #[clap(short, long)]
    pub task: String,

    /// Module is the name of the dagger module to generate.
    #[clap(short, long)]
    pub module: Option<String>,

    /// Module type is the type of the module to generate.
    #[clap(long, default_value = "full")]
    pub module_type: String,
}
