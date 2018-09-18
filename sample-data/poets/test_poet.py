import argparse
import json

parser = argparse.ArgumentParser()
group = parser.add_mutually_exclusive_group()
group.add_argument("--write", action="store_true")
group.add_argument("--critique", type=str, help="rate a poem between 0-1")
group.add_argument("--study", type=str, help="learn from new poems")
args = parser.parse_args()


def write_poem():
    poem = {
        'title': 'buffalo buffalo',
        'content': "buffalo buffalo buffalo buffalo buffalo buffalo",
    }

    s = json.dumps(poem)

    print(s)


def critique_poem(poem):
    critique = {
        'score': 0.56,
    }

    s = json.dumps(critique)

    print(s)


def study_poem():
    update = {
        'success': True,
    }

    s = json.dumps(update)

    print(s)


if args.write:
    write_poem()
elif args.critique:
    critique_poem(args.critique)
elif args.study:
    study_poem()
