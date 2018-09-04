import argparse

parser = argparse.ArgumentParser()
group = parser.add_mutually_exclusive_group()
group.add_argument("--write", action="store_true")
group.add_argument("--critique", type=str, help="rate a poem between 0-1")
group.add_argument("--study", type=str, help="learn from new poems")
args = parser.parse_args()



print(args)
