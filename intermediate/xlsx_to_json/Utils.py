import os
import subprocess
from distutils.log import warn
import re


def write_file(out_path: str, content: str):
    with open(out_path, 'w', encoding='utf-8') as file:
        file.write(content)


def write_file_golang_struct(file_name, data: list, out_ath: str, package_name: str):
    if not os.path.exists(out_ath):
        os.makedirs(out_ath)
    file_path = out_ath + '/' + file_name + '.go'  # 设置路径
    content: str = 'package {} \n\n'.format(package_name)
    for item in data:
        content += item
    write_file(file_path, content)


def write_file_suffix(file_name, content: str, out_ath: str, suffix: str):
    if not os.path.exists(out_ath):
        os.makedirs(out_ath)
    file_path = out_ath + '/' + file_name + '.' + suffix  # 设置路径
    write_file(file_path, content)


# 复制文件夹或者目录
def copy(path_1, path_2):
    a, b = subprocess.getstatusoutput('cp -R -f ' + path_1 + " " + path_2)
    if (a != 0):
        print("\033[0;31m" + "[error]: " + b + "\033[0m")
        return False
    return True


# 删除某个路径下的所有文件
def remove_all_file(filepath):
    del_list = os.listdir(filepath)
    for f in del_list:
        file_path = os.path.join(filepath, f)
        if os.path.isfile(file_path): os.remove(file_path)


# 补齐字段长度
def __fill_up(str: str):
    length = len(str)
    num = 72 - length
    for i in range(0, num): str += " "
    return str


# 返回字节大小
def get_file_size(file_path):
    fsize = os.path.getsize(file_path)  # 返回的是字节大小
    if fsize < 1024:
        return (round(fsize, 2), 'byte')
    kb_value = fsize / 1024
    if kb_value < 1024:
        return (round(kb_value, 2), 'kb')
    mb_value = kb_value / 1024
    if mb_value < 1024:
        return (round(mb_value, 2), 'mb')
    return (round(mb_value / 1024), 'gb')


# 判断Json大小并打印信息
def larger_size_and_print(_size, _file):
    log_str = "parsing file successful: " + _file
    log_str = __fill_up(log_str)
    log_str = log_str + "\t json size: " + str(_size[0]) + " " + _size[1]
    if (_size[1] != "byte" and _size[0] > 500):
        warn("\033[0;31m" + log_str + "\033[0m")
    else:
        print(log_str)


# 获取某个路径下的所有xlsx文件
def get_all_file_path(root: str, start: str = '~$,.~', end: str = 'xlsx'):
    name = []
    path = []
    for file_path, dir_names, file_names in os.walk(root):
        for filename in file_names:
            for pre in start.split(","):
                if str.startswith(filename, pre): continue
            if str.startswith(filename, pre): continue
            if str.endswith(filename, end):
                name.append(filename)
                path.append(os.path.join(file_path, filename))
    return [name, path]


# 首字母小写
def decapitalize(str: str): return str[:1].lower() + str[1:]


# 返回注释
def get_ts_annotation(file_name, num: 0):
    str = ""
    for i in range(0, num): str += "\t"
    return str + "/**\n" + str + " * " + file_name + "\n" + str + " */"


def snakecase(string):
    """Convert string into snake case.
    Join punctuation with underscore

    Args:
        string: String to convert.

    Returns:
        string: Snake cased string.

    """

    string = re.sub(r"[\-\.\s]", '_', str(string))
    if not string:
        return string
    return lowercase(string[0]) + re.sub(r"[A-Z]", lambda matched: '_' + lowercase(matched.group(0)), string[1:])


def lowercase(string):
    """Convert string into lower case.

    Args:
        string: String to convert.

    Returns:
        string: Lowercase case string.

    """

    return str(string).lower()


def pascalcase(string):
    """Convert string into pascal case.

    Args:
        string: String to convert.

    Returns:
        string: Pascal case string.

    """

    return capitalcase(camelcase(string))


def camelcase(string):
    """ Convert string into camel case.

    Args:
        string: String to convert.

    Returns:
        string: Camel case string.

    """

    string = re.sub(r"\w[\s\W]+\w", '', str(string))
    if not string:
        return string
    return lowercase(string[0]) + re.sub(r"[\-_\.\s]([a-z])", lambda matched: uppercase(matched.group(1)), string[1:])


def uppercase(string):
    """Convert string into upper case.

    Args:
        string: String to convert.

    Returns:
        string: Uppercase case string.

    """

    return str(string).upper()


def capitalcase(string):
    """Convert string into capital case.
    First letters will be uppercase.

    Args:
        string: String to convert.

    Returns:
        string: Capital case string.

    """

    string = str(string)
    if not string:
        return string
    return uppercase(string[0]) + string[1:]
