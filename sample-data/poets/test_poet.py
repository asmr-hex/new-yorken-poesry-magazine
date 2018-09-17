import argparse

parser = argparse.ArgumentParser()
group = parser.add_mutually_exclusive_group()
group.add_argument("--write", action="store_true")
group.add_argument("--critique", type=str, help="rate a poem between 0-1")
group.add_argument("--study", type=str, help="learn from new poems")
args = parser.parse_args()


def write_poem():
    print("buffalo buffalo buffalo buffalo buffalo buffalo")


def critique_poem(poem):
    print(0.33)


def study_poem():
    print("true")


if args.write:
    write_poem()
elif args.critique:
    critique_poem(args.critique)
elif args.study:
    study_poem()
