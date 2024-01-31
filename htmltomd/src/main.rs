use std::{env,fs};
use reqwest;
use html2md;
fn main() {
    // 获取命令行参数
    let args: Vec<String> = env::args().collect();
    // 获取网址
    let url: &str = &args[1];
    // 获取输出文件名
    let output: &str = &args[2];
    let body = reqwest::blocking::get(url)
    .unwrap().text().unwrap();
    let md = html2md::parse_html(&body);
    fs::write(output, md.as_bytes()).unwrap();
}
