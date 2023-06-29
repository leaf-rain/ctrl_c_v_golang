import re

import json

import math
import Config
import Utils
import sys
import getopt
import yaml
import pandas as pd

Table = 'table'
Object = 'object'
KV = 'kv'

ConfigType = [Table, Object, KV]


# 根据配置Json文件初始化配置
def init_config(config_path: str):
    if config_path == '': return
    Config.config_file_path = config_path
    with open(config_path, 'r+', encoding='utf-8') as file:
        data = yaml.load(file, Loader=yaml.FullLoader)
        if "in_path_xlsx" in data: Config.in_path_xlsx = data["in_path_xlsx"]
        if "out_put_path" in data: Config.out_put_path = data["out_put_path"]
        if "out_opt_type" in data: Config.out_opt_type = data["out_opt_type"]
        if "ignore_excel_file" in data: Config.ignore_excel_file = data["ignore_excel_file"]
        if "golang_package_name" in data: Config.golang_package_name = data["golang_package_name"]
        if "columns" in data: Config.columns = set(data["columns"])


def main():
    # 声明变量
    config_path: str = './local_config.yaml'
    # 加载命令行参数
    opts, _ = getopt.getopt(sys.argv[1:], "hc:", ["config="])  # 长选项模式
    for opt, arg in opts:
        if opt == '-h':
            print('Main.py -c <config file>')
            sys.exit()
        elif opt in ("-c", "--config"):
            config_path = arg
    # 开始运行
    print("-------------加载环境配置文件:{}--------------------".format(config_path))
    init_config(config_path)
    for item in Config.in_path_xlsx:
        file_and_path = Utils.get_all_file_path(item)
        if len(file_and_path) == 0 or len(file_and_path[0]) == 0 or len(file_and_path[0]) != len(file_and_path[1]):
            print("-------------文件路径加载有问题:{}--------------------".format(config_path))
            continue
        for index in range(0, len(file_and_path[0])):
            file_name = file_and_path[0][index]
            file_path = file_and_path[1][index]
            config_name = file_name.split(".")[0]
            if file_name in Config.ignore_excel_file:  # 跳过不需要过滤文件
                print("-------------跳过xlsx文件:{}--------------------".format(file_and_path[0][index]))
                continue
                # 读取xlsx文件
            print("-------------读取xlsx文件:{}--------------------".format(file_and_path[0][index]))
            data = pd.read_excel(file_path, header=0)
            config_list = []
            if data.columns[0] in ConfigType:
                one_config = ClassInfo()
                one_config.type = data.columns[0]
                one_config.name = data.columns[1]
                one_config.remark = data.columns[2]
                one_config.user = set(str.split(data.columns[3], ","))
                one_config.start = 0
                config_list.append(one_config)
            for data_index in range(0, len(data.iloc[:, [0]].values)):
                tmp_value = data.iloc[:, [0]].values[data_index][0]
                if len(config_list) > 0:
                    config_list[-1].end = data_index
                if isinstance(tmp_value, str) and (tmp_value in ConfigType):
                    new_config = ClassInfo()
                    new_config.type = data.iloc[:, [0]].values[data_index][0]
                    new_config.name = data.iloc[:, [1]].values[data_index][0]
                    new_config.remark = data.iloc[:, [2]].values[data_index][0]
                    new_config.start = data_index + 1  # 需要跳过标识列
                    user_param = data.iloc[:, [3]].values[data_index][0]
                    if not (isinstance(user_param, float) and math.isnan(user_param)):
                        new_config.user = set(str.split(data.iloc[:, [3]].values[data_index][0], ","))
                    config_list.append(new_config)
            config_list[-1].end += 1  # 最后一列需要再+1
            # 切割pandas数据
            for config_list_index in range(0, len(config_list)):
                if config_list[config_list_index].type == Table:
                    config_list[config_list_index].content = data.iloc[config_list[config_list_index].start:config_list[
                        config_list_index].end, :]
                    config_list[config_list_index].content.columns = config_list[0].content.iloc[0, :].values.tolist()
                    config_list[config_list_index].content.index = range(0,
                                                                         config_list[config_list_index].content.shape[
                                                                             0])
                elif config_list[config_list_index].type == KV or config_list[config_list_index].type == Object:
                    config_list[config_list_index].content = data.iloc[config_list[config_list_index].start:config_list[
                        config_list_index].end, :4]
                    config_list[config_list_index].content.columns = ["key", "value", "type", "user"]
                    config_list[config_list_index].content.index = range(0,
                                                                         config_list[config_list_index].content.shape[
                                                                             0])
            # 解析成map树
            config_map, out_put_name = parse_map(config_list)
            if out_put_name == '':
                out_put_name = config_name
            for out_index in range(0, len(Config.out_opt_type)):
                if out_index < len(Config.out_put_path):
                    out_path = Config.out_put_path[out_index]
                else:
                    out_path = Config.out_put_path[0]
                if Config.out_opt_type[out_index] == 1:  # json文件
                    js = json.dumps(config_map, ensure_ascii=False)
                    Utils.write_file_suffix(out_put_name, js, out_path, "json")
                elif Config.out_opt_type[out_index] == 2:  # yaml文件
                    js = yaml.dump(config_map, allow_unicode=True)
                    Utils.write_file_suffix(out_put_name, js, out_path, "yaml")
            print("-------------读取xlsx文件success:{}--------------------".format(file_and_path[0][index]))


def parse_map(config_list: list):
    result: dict = {}
    out_put_name: str = ''
    for item in config_list:
        if len(item.user) > 0 and len(Config.columns) > 0:
            if len(set(item.user).intersection(set(Config.columns))) == 0:
                continue
        if item.type == Table:
            if not isinstance(item.content, pd.DataFrame):
                print(
                    "[parse_map] Table item content type failed, item name:{}, item type:{} item content type:{}".format(
                        item.name, item.type, type(item.content)))
                continue
            # 解析备注
            remark = json.loads(item.remark)
            type_index: int = 0
            user_index: int = 0
            remarks_index: str = ''
            data_start_index: int = 0
            if "type" in remark: type_index = int(remark["type"])
            if "user" in remark: user_index = int(remark["user"])
            if "remarks" in remark: remarks_index = str(remark["remarks"])
            if "data_start" in remark: data_start_index = int(remark["data_start"])
            if "out_put_name" in remark: out_put_name = str(remark["out_put_name"])
            max_j = item.content.shape[0]
            max_i = item.content.shape[1]
            if type_index == 0 or data_start_index == 0 or max_j <= 0 or max_i <= 0 or data_start_index > max_j:
                print(
                    "[parse_map] table parse failed,item name:{}, item type:{}, type_index:{}, data_start_index:{}, max_i:{}, max_j:{}, data_start_index:{}".format(
                        item.name, item.type, type_index, data_start_index, max_i, max_j, data_start_index))
                continue
            result[item.name]: list = []
            for j in range(data_start_index, max_j):
                data: dict = {}
                for i in range(0, max_i):
                    if user_index > 0 :
                        user_list = parse_column(item.content.at[user_index, item.content.columns[i]], "array")
                        if len(user_list) > 0 and len(Config.columns) > 0:
                            if len(set(user_list).intersection(set(Config.columns))) == 0:
                                continue
                    tmpValue = parse_column(item.content.at[j, item.content.columns[i]],
                                                                 item.content.at[type_index, item.content.columns[i]])
                    if not tmpValue is None:
                        data[item.content.columns[i]] = tmpValue
                result[item.name].append(data)
        elif item.type == KV:
            if not isinstance(item.content, pd.DataFrame):
                print("[parse_map] KV item content type failed, item name:{}, item type:{} item content type:{}".format(
                    item.name, item.type, type(item.content)))
                continue
            max_j = item.content.shape[0]
            for j in range(0, max_j):
                user_list = parse_column(item.content.at[j, item.content.columns[3]], "array")
                if len(user_list) > 0 and len(Config.columns) > 0:
                    if len(set(user_list).intersection(set(Config.columns))) == 0:
                        continue
                tmpValue = parse_column(
                    item.content.at[j, item.content.columns[1]], item.content.at[j, item.content.columns[2]])
                if not tmpValue is None:
                    result[parse_column(item.content.at[j, item.content.columns[0]], "str")] = tmpValue
        elif item.type == Object:
            result[item.name]: dict = {}
            if not isinstance(item.content, pd.DataFrame):
                print(
                    "[parse_map] Object item content type failed, item name:{}, item type:{} item content type:{}".format(
                        item.name, item.type, type(item.content)))
                continue
            max_j = item.content.shape[0]
            for j in range(0, max_j):
                user_list = parse_column(item.content.at[j, item.content.columns[3]], "array")
                if len(user_list) > 0 and len(Config.columns) > 0:
                    if len(set(user_list).intersection(set(Config.columns))) == 0:
                        continue
                tmpValue = parse_column(
                    item.content.at[j, item.content.columns[1]], item.content.at[j, item.content.columns[2]])
                if not tmpValue is None:
                    result[item.name][parse_column(item.content.at[j, item.content.columns[0]], "str")] = tmpValue
    return result, out_put_name


def parse_column(data, ty: str):
    try:
        if ty == "number" or ty == "num" or ty == "int":
            if isinstance(data, float) and math.isnan(data): return 0
            return int(data)
        if ty == "float":
            if isinstance(data, float) and math.isnan(data): return 0.0
            return float(data)
        if ty == "string" or ty == "str":
            if isinstance(data, float) and math.isnan(data): return ''
            return str(data)
        elif ty == "bool":
            if isinstance(data, float) and math.isnan(data): return False
            return bool(data)
        elif ty == "json" or ty == "js":
            if isinstance(data, float) and math.isnan(data): return None
            data = re.sub(r'\s', '', data)
            return json.loads(data)
        elif ty == "array" or ty == "list":
            if isinstance(data, float) and math.isnan(data): return []
            return str.split(str(data), ",")
    except Exception as err:
        print(
            "[parse_column] err:{} 类型转换错误, 输入type:{}, 输入值:{}, 输出type:{}".format(err, type(data), data, ty))
        if ty == "number" or ty == "num" or ty == "int":
            return 0
        if ty == "float":
            return 0.0
        if ty == "string" or ty == "str":
            return ''
        elif ty == "bool":
            return False
        elif ty == "json" or ty == "js":
            return {}
        elif ty == "array" or ty == "list":
            return []


class ClassInfo:
    def __int__(self):
        self.name = ""
        self.type = ""
        self.start = 0
        self.end = 0
        self.remark = ""
        self.user = set()
        self.content = pd.DataFrame()

    user = set()

    def __str__(self):
        return "name:{}, type:{}, start:{}, end:{}".format(self.name, self.type, self.start, self.end)


# def load_json(config_path: str):

# 开始运行
if __name__ == '__main__':
    main()
    # print(parse_column("abc", "int"))
