# python3
import github    # needs https://pygithub.readthedocs.io/en/latest/introduction.html#download-and-install
import itertools
import operator
import re
import subprocess
from typing import Mapping, Sequence, Union

REPO = "tools"
# "googlecodelabs"
REPO_OWNER = "Marcial1234"
GITHUB_TOKEN = "dec8f5906be7a096e3c2f548b525aad8e0f0aa97"

def main():
  # Create a Github using an access token
  g = github.Github(GITHUB_TOKEN)
  repo = g.get_repo("{}/{}".format(REPO_OWNER, REPO))

  # chains = { pr-chain-name: {pr#: [ref, branch]} }
  chains = {}
  chained_prs_regex = re.compile(r"\[(\w+(?:-\w+)+) - Part(\s?\d+)\]")
  
  for pr in repo.get_pulls(state="closed", sort='created', direction="asc"):
    match = chained_prs_regex.match(pr.title)
    if match:
      add_to_chain(pr, chains, *match.groups())

  for pr_chain, g in chains.items():
    branches_to_be_deleted, bases = get_branches_from_labels(g.values())
    # repo.get_branches <= see if these branches exist on the first place...
    # pop_already_deleted(chain)
    if is_sequential(g.keys()):
      print("'{}' branch chain can be deleted!".format(pr_chain))

      print_branches_to_delete(branches_to_be_deleted)
      if input("Proceed? y/n: ") == 'y':
        if type(branches_to_be_deleted) == str:
          delete_branch(branches_to_be_deleted)
        else:
          delete_chain(branches_to_be_deleted, bases)
      else:
        print("Cancelled! Pheww!")
    else:
      print("'{}' chain cannot be deleted yet!!".format(pr_chain))
      print("Sequence: {}".format(list(g.items())))


def add_to_chain(
  pr, tracking_chain: Mapping[str, Mapping[int, Sequence[str]]],
  new_chain_name: str, number: int) -> None:
  data = [pr.head.label, pr.base.label]

  if new_chain_name not in tracking_chain:
    tracking_chain[new_chain_name] = {int(number): data}
  else:
    tracking_chain[new_chain_name][int(number)] = data


def is_sequential(data: Sequence[int]) -> bool:
  asc_data = list(data)
  asc_data.reverse()
  seq = [g for key, g in itertools.groupby(enumerate(asc_data), difference)]
  return len(seq) == 1


def difference(args: Sequence[Union[int, str]]) -> int:
  index, n = args
  return index - int(n)


def get_branches_from_labels(labels: Mapping[int, Sequence[str]]):
  labels = list(labels)
  labels.sort(reverse=True)
  if len(labels) == 1: return labels[0]
  return operator.itemgetter(1)(labels), operator.itemgetter(0)(labels)


def print_branches_to_delete(branches: Sequence[str]):
  if type(branches) == str: branches = [branches]
  branch_list = "\n\t".join('"{0}"'.format(b) for b in branches)
  print("Branches to delete:\n\t{}\n".format(branch_list))


def delete_chain(
  parent_branches: Sequence[str], reverse_sequential_branches: Sequence[str]):
  # last check
  if parent_branches[:-1] == reverse_sequential_branches[1:]:
    for b in reverse_sequential_branches:
      print()
      delete_branch(b)
      print()
  else:
    print("nope!")

def delete_branch(raw_branch: str):
  # works as either "label" (user:branch) or "ref" (branch)
  branch = raw_branch.split(":")[-1]
  cmd = ["git", "push", "--delete", "origin", branch]
  print("Deleting '{}'".format(raw_branch))
  subprocess.check_output(cmd, shell=True)

if __name__ == '__main__':
  main()