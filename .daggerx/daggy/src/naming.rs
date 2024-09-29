pub fn capitalize_module_name(module_name: &str) -> String {
    let mut chars = module_name.chars();
    match chars.next() {
        None => String::new(),
        Some(first) => first.to_uppercase().collect::<String>() + chars.as_str(),
    }
}

pub fn to_camel_case(s: &str) -> String {
    s.split('-')
        .enumerate()
        .map(|(i, part)| {
            if i == 0 {
                part.to_lowercase()
            } else {
                capitalize_module_name(part)
            }
        })
        .collect()
}

pub fn to_pascal_case(s: &str) -> String {
    s.split('-').map(capitalize_module_name).collect()
}
