import os
from translate import Translator
import json



def translatei18n(translator, to="en"):
    folder = "/Users/wsh/Public/markless-github/local/translations/"
    zh_dic = json.load(open(os.path.join(folder, "zh-CN.json"), 'r', encoding='utf-8'))

    existed_file = os.path.join(folder, f"{to}.json")
    if os.path.exists(existed_file):
        existed_dic = json.load(open(existed_file, 'r', encoding='utf-8'))
    else:
        existed_dic = {}

    new_dic = {}

    for k, v in zh_dic.items():
        print(k + '：' + v)
        if existed_dic.get(k):
            new_dic[k] = existed_dic[k]
            continue
        else:
            try:
                val = translator.translate(v.replace(' ', ''))
                new_dic[k] = val
                print(f'translate:{k}：{val}')
            except Exception as e:
                new_dic[k] = v

    with open(os.path.join(folder, f"{to}_temp.json"), 'w', encoding='utf-8') as f:
        json.dump(new_dic, f, ensure_ascii=False, indent=4)



# translator = Translator(from_lang="Chinese",to_lang="English")
translator = Translator(from_lang="Chinese",to_lang="Japanese")
translatei18n(translator, to="ja")
