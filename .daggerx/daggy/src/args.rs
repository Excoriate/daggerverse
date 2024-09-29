use clap::Parser;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
pub struct Args {
    /// Task is the name of the task to run
    #[arg(short = 't', long = "task")]
    pub task: String,

    /// Module is the name of the dagger module to generate.
    #[arg(short = 'm', long = "module")]
    pub module: Option<String>,

    /// Module type is the type of the module to generate.
    #[arg(short = 'y', long = "type")]
    pub module_type: Option<String>,
}
