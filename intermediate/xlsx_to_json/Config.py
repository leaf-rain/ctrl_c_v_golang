# 配置文件位置
config_file_path = "./config.yaml"

# xlsx文件path
in_path_xlsx = ["../"]

# 输出路径, 列表，与下方输出文件匹配的话则对应各个输出位置，否则默认使用第一个输出位置
out_put_path = ["../"]

# 输出文件, 1:json, 2:yaml, 3:golang struct
out_opt_type = [1, 2, 3]

# 字段选择
columns = set()

# 解析时忽略的文件
ignore_excel_file = []

# 解析后golang的package名称
golang_package_name = "consts"
