#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import argparse
import json
import re
import time
import subprocess
import os
from collections import Counter
from pathlib import Path
from typing import Dict, Any, List, Tuple

import logging
from datetime import datetime

# 配置错误日志
error_log_path = Path(f"translate_errors_{datetime.now()}.log")
logging.basicConfig(
    level=logging.ERROR,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler(error_log_path, encoding='utf-8'),
        logging.StreamHandler()  # 同时输出到控制台
    ]
)
error_logger = logging.getLogger(__name__)

# 支持的语言映射
SUPPORTED_LANGUAGES = {
    'zh-cn': '中文（简体）',
    'zh-tw': '中文（繁体）',
    'en': '英语',
    'es': '西班牙语',
    'fr': '法语',
    'de': '德语',
    'ja': '日语',
    'ko': '韩语',
    'ru': '俄语',
    'pt': '葡萄牙语',
    'it': '意大利语',
    'ar': '阿拉伯语',
    'hi': '印地语',
    'th': '泰语',
    'vi': '越南语',
}

from deep_translator import GoogleTranslator
translation_cache = {}

def flatten_dict(data: Dict[str, Any], parent_key: str = '', sep: str = '.') -> Dict[str, str]:
    """
    扁平化嵌套字典
    例如：{'common': {'name': '中文'}} -> {'common.name': '中文'}
    """
    items = []
    for k, v in data.items():
        new_key = f"{parent_key}{sep}{k}" if parent_key else k
        if isinstance(v, dict):
            items.extend(flatten_dict(v, new_key, sep=sep).items())
        elif isinstance(v, str):
            items.append((new_key, v))
        else:
            items.append((new_key, str(v)))
    return dict(items)

def translate_text(text: str, target_lang: str = 'es') -> str:
    """翻译文本，带缓存和重试"""
    if not text or not re.search(r'[一-龥]', text):
        return text
    
    cache_key = f"{text}_{target_lang}"
    if cache_key in translation_cache:
        return translation_cache[cache_key]
    
    max_retries = 3
    for attempt in range(max_retries):
        try:
            translated = GoogleTranslator(source='zh-CN', target=target_lang).translate(text)
            print(f"中文: {text}")
            print(f"翻译: {translated}")
            translation_cache[cache_key] = translated
            time.sleep(0.1)
            return translated
        except Exception as e:
            print(f"翻译失败 (尝试 {attempt+1}/{max_retries}): {text[:50]}... 错误: {e}")
            time.sleep(1)
    
    print(f"翻译失败，保留原文: {text[:50]}...")
    return text

def process_php_file(file_path: Path, target_lang: str) -> Dict[str, Any]:
    """处理单个PHP文件，提取并翻译"""
    print(f"处理PHP文件: {file_path}")
    
    file_abs_path = str(file_path.absolute())
    php_command = (
        f'php -r "echo json_encode(require \'{file_abs_path}\', JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);"'
    )
    result = subprocess.run(
        php_command,
        capture_output=True,
        text=True,
        encoding='utf-8',
        shell=True,
        check=True
    )
    json_data = json.loads(result.stdout)
     # 扁平化嵌套字典
    json_data = flatten_dict(json_data)
    translations = {}
    for key, chinese_text in json_data.items():
        translated_text = translate_text(chinese_text, target_lang)
        translations[key] = {
            'original': chinese_text,
            'translated': translated_text
        }

    print(f"  完成: {len(json_data)} 个字符串")
    
    return {
        'file_path': str(file_path),
        'translations': translations,
        'total_strings': len(json_data)
    }

def process_json_file(file_path: Path, target_lang: str) -> Dict[str, Any]:
    """处理JSON文件，全部重新翻译"""
    print(f"处理JSON文件: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        json_data = json.load(f)
    
    translations = {}
    for key, value in json_data.items():
        original = value.get('original', '') if isinstance(value, dict) else value
        translated_text = translate_text(original, target_lang)
        
        translations[key] = {
            'original': original,
            'translated': translated_text
        }
    
    print(f"  完成: {len(json_data)} 个字符串")
    
    return {
        'file_path': str(file_path),
        'translations': translations,
        'total_strings': len(json_data)
    }

def collect_php_files(paths: List[Path]) -> List[Path]:
    """收集PHP文件"""
    files = []
    for path in paths:
        if not path.exists():
            print(f"警告: 路径不存在 - {path}")
            continue
        if path.is_file() and path.suffix == '.php':
            files.append(path)
        elif path.is_dir():
            for php_file in path.rglob("*.php"):
                if not any(part.startswith(".") for part in php_file.relative_to(path).parts):
                    files.append(php_file)
        else:
            print(f"警告: 跳过非PHP文件 - {path}")
    return sorted(set(files))

def collect_json_files(paths: List[Path]) -> List[Path]:
    """收集JSON文件"""
    files = []
    for path in paths:
        if not path.exists():
            print(f"警告: 路径不存在 - {path}")
            continue
        if path.is_file() and path.suffix == '.json':
            files.append(path)
        elif path.is_dir():
            for json_file in path.rglob("*.json"):
                if not any(part.startswith(".") for part in json_file.relative_to(path).parts):
                    files.append(json_file)
        else:
            print(f"警告: 跳过非JSON文件 - {path}")
    return sorted(set(files))

def get_file_key(file_path: Path, base_path: Path = None) -> str:
    """获取文件的唯一标识键"""
    if base_path:
        try:
            rel_path = file_path.relative_to(base_path)
            return str(rel_path.with_suffix('')).replace('\\', '/')
        except ValueError:
            pass
    return file_path.stem

def build_output_data(files: List[Path], results: List[Dict], base_path: Path = None, mode: str = "php") -> Dict:
    """构建最终的JSON输出数据"""
    output_data = {}
    
    for file_path, result in zip(files, results):
        if not result or result['total_strings'] == 0:
            continue
        
        file_key = get_file_key(file_path, base_path)
        
        for key, translation in result['translations'].items():
            if mode == "json":
                full_key = key
            else:
                full_key = f"{file_key}.{key}"
            
            output_data[full_key] = {
                'file': str(file_path),
                'file_key': file_key,
                'key': key,
                'original': translation['original'],
                'translated': translation['translated']
            }
    
    return output_data

def print_summary(files: List[Path], mode: str) -> None:
    """打印文件统计"""
    if not files:
        print(f"没有找到要处理的{mode.upper()}文件")
        return
    
    directories = Counter([f.parent for f in files])
    print(f"找到 {len(files)} 个{mode.upper()}文件")
    for dir_path, count in sorted(directories.items()):
        print(f"  {dir_path}: {count} 个文件")

def main():
    parser = argparse.ArgumentParser(
        description="将中文翻译成指定语言，支持PHP语言包或JSON文件",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    
    parser.add_argument(
        "--path", "-p", 
        action="append", 
        dest="paths",
        help="指定要处理的文件或文件夹路径（可多次使用）"
    )
    parser.add_argument(
        "--lang", "-l", 
        type=str, 
        default="es",
        choices=SUPPORTED_LANGUAGES.keys(),
        help="目标语言代码（默认: es 西班牙语）"
    )
    parser.add_argument(
        "--output", "-o", 
        type=str, 
        default="translations.json",
        help="输出JSON文件路径（默认: translations.json）"
    )
    parser.add_argument(
        "--mode", "-m",
        type=str,
        default="php",
        choices=["php", "json"],
        help="输入模式：php（处理PHP文件）或 json（处理JSON文件）（默认: php）"
    )
    
    args = parser.parse_args()
    
    target_lang_name = SUPPORTED_LANGUAGES.get(args.lang, args.lang)
    print(f"目标语言: {args.lang} ({target_lang_name})")
    print(f"运行模式: {args.mode.upper()}")
    
    # 确定要处理的路径
    if args.paths:
        target_paths = []
        for p in args.paths:
            path = Path(p)
            if not path.is_absolute():
                path = Path.cwd() / path
            target_paths.append(path)
    else:
        if args.mode == "php":
            default_path = Path.cwd() / "zh"
        else:
            default_path = Path.cwd() / "translations.json"
        
        if not default_path.exists():
            print(f"错误: 默认路径不存在 - {default_path}")
            print("请使用 --path 参数指定要处理的文件或文件夹")
            return
        target_paths = [default_path]
    
    # 根据模式收集文件
    if args.mode == "php":
        files = collect_php_files(target_paths)
        process_func = process_php_file
    else:
        files = collect_json_files(target_paths)
        process_func = process_json_file
    
    if not files:
        print(f"没有找到任何{args.mode.upper()}文件")
        return
    
    print_summary(files, args.mode)
    
    # 确定基准目录（仅PHP模式）
    base_dir = None
    if args.mode == "php":
        current_dir = os.getcwd()
        # 拼接上 "zh"
        base_dir = os.path.join(current_dir, "zh")      
    
    print("\n开始提取并翻译...")
    print("-" * 50)
    
    results = []
    total_strings = 0
    success_count = 0
    error_count = 0
    
    for i, file_path in enumerate(files, 1):
        try:
            result = process_func(file_path, args.lang)
            if result:
                results.append(result)
                total_strings += result['total_strings']
                success_count += 1
            else:
                results.append(None)
                error_count += 1
            
            if i % 10 == 0 or i == len(files):
                print(f"进度: {i}/{len(files)} 文件已处理, 已处理 {total_strings} 个字符串")
            
        except Exception as e:
            print(f"处理文件失败 {file_path}: {e}")
            error_logger.error(f"处理文件失败: {file_path}")
            results.append(None)
            error_count += 1
        
        if i % 5 == 0:
            time.sleep(0.5)
    
    # 构建输出数据
    output_data = build_output_data(files, results, base_dir, args.mode)
    
    # 保存到JSON文件
    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(
        json.dumps(output_data, ensure_ascii=False, indent=2),
        encoding='utf-8'
    )
    
    print("\n" + "=" * 50)
    print("处理完成!")
    print(f"成功处理: {success_count} 个文件")
    print(f"处理失败: {error_count} 个文件")
    print(f"翻译字符串总数: {total_strings}")
    print(f"输出文件: {output_path.absolute()}")
    
    if output_data:
        files_count = len(set([v['file'] for v in output_data.values()]))
        print(f"涉及文件数: {files_count}")

if __name__ == "__main__":
    main()