# SPDX-License-Identifier: Apache-2.0

"""
FaultPlane Traction Hunter

Discovers public, technically relevant network/reliability pain points from:
  - GitHub Issues
  - GitHub Discussions
  - GitHub Wiki

The tool:
  * scores relevance
  * extracts public GitHub usernames associated with relevant discussions
  * generates human-review outreach drafts
  * writes structured JSON/Markdown artifacts

It deliberately DOES NOT:
  * automatically post comments
  * automatically send DMs
  * automatically mention users
  * automatically open issues
  * automatically spam repositories
"""

from __future__ import annotations

import json
import os
import re
import time
from collections import defaultdict
from pathlib import Path
from typing import Any

import requests
from github import Github
from github.GithubException import GithubException


REPOSITORY = "devloperdevesh/FaultPlane"

TARGET_REPOSITORIES = [
    "langchain-ai/langchain",
    "crewAIInc/crewAI",
    "embedchain/embedchain",
    "tiangolo/fastapi",
    "microsoft/autogen",
    "run-llama/llama_index",
]

FAILURE_SIGNALS = {
    "timeout": 5,
    "connection timeout": 6,
    "connection reset": 6,
    "connection drop": 6,
    "connection refused": 5,
    "socket error": 6,
    "broken pipe": 5,
    "read timeout": 5,
    "write timeout": 5,
    "network error": 4,
    "connection pool": 4,
    "pool exhausted": 6,
    "too many connections": 6,
    "transport error": 4,
    "upstream error": 4,
    "http 502": 4,
    "http 503": 4,
    "http 504": 5,
    "eof": 3,
    "agent crash": 4,
}

INFRASTRUCTURE_SIGNALS = {
    "async": 2,
    "concurrent": 2,
    "concurrency": 2,
    "streaming": 2,
    "worker": 1,
    "workers": 1,
    "tcp": 2,
    "socket": 3,
    "sockmap": 4,
    "ebpf": 4,
    "linux": 2,
    "kubernetes": 2,
    "cluster": 2,
    "load": 2,
    "high traffic": 3,
    "production": 3,
    "scale": 2,
    "scaling": 2,
    "retry": 2,
    "retries": 2,
}

MAX_ISSUES_PER_REPOSITORY = 8
MAX_DISCUSSIONS_PER_REPOSITORY = 8
MAX_TOTAL_RESULTS = 100

OUTPUT_DIR = Path("artifacts/traction")
JSON_OUTPUT = OUTPUT_DIR / "traction_matches.json"
MARKDOWN_OUTPUT = OUTPUT_DIR / "traction_report.md"

GRAPHQL_URL = "https://api.github.com/graphql"


def normalize(value: str | None) -> str:
    return re.sub(r"\s+", " ", value or "").strip()


def score_text(title: str, body: str) -> tuple[int, list[str]]:
    title_l = normalize(title).lower()
    body_l = normalize(body).lower()
    combined = f"{title_l} {body_l}"

    score = 0
    signals: list[str] = []

    for signal, weight in FAILURE_SIGNALS.items():
        if signal in combined:
            signals.append(signal)
            score += weight

            if signal in title_l:
                score += 2

    for signal, weight in INFRASTRUCTURE_SIGNALS.items():
        if signal in combined:
            score += weight

    return score, signals


def build_draft(
    repository: str,
    number: int | str,
    title: str,
    url: str,
    signals: list[str],
) -> str:
    signal_text = ", ".join(signals[:6])

    return f"""Hey @username,

We came across this while investigating real-world reliability problems
in AI and distributed application workloads.

The failure signals that stood out here are: {signal_text}.

We're building FaultPlane, an open-source Go/Linux networking project focused
on failure handling and eBPF-based socket/data-path integration.

The goal is to move selected failure-handling paths closer to the Linux
networking boundary instead of relying entirely on application-level retry
loops.

We're currently measuring the implementation across different workloads,
so we don't want to claim a guaranteed latency, uptime, or allocation result
for this specific environment without reproducing it first.

If this problem is still reproducible, we'd be interested in comparing the
failure mode with FaultPlane and seeing whether it maps to the class of
network failure we're targeting.

FaultPlane:
https://github.com/{REPOSITORY}

Related issue:
{url}

Signals detected:
{signal_text}

If useful, we can help map the failure path and identify where the connection
is being dropped, reset, or timing out.

— FaultPlane
"""


def github_headers(token: str) -> dict[str, str]:
    return {
        "Authorization": f"Bearer {token}",
        "Accept": "application/vnd.github+json",
        "X-GitHub-Api-Version": "2022-11-28",
    }


def graphql(
    token: str,
    query: str,
    variables: dict[str, Any],
) -> dict[str, Any]:
    response = requests.post(
        GRAPHQL_URL,
        headers=github_headers(token),
        json={
            "query": query,
            "variables": variables,
        },
        timeout=30,
    )

    response.raise_for_status()

    payload = response.json()

    if payload.get("errors"):
        raise RuntimeError(payload["errors"])

    return payload["data"]


def collect_discussions(
    token: str,
    owner: str,
    name: str,
) -> list[dict[str, Any]]:
    query = """
    query($owner: String!, $name: String!) {
      repository(owner: $owner, name: $name) {
        discussions(first: 20, orderBy: {field: UPDATED_AT, direction: DESC}) {
          nodes {
            number
            title
            body
            url
            updatedAt
            author {
              login
            }
            category {
              name
            }
            comments(first: 20) {
              nodes {
                body
                author {
                  login
                }
              }
            }
          }
        }
      }
    }
    """

    data = graphql(
        token,
        query,
        {
            "owner": owner,
            "name": name,
        },
    )

    repository = data.get("repository")

    if not repository:
        return []

    discussions = repository.get("discussions", {}).get("nodes", [])

    results = []

    for discussion in discussions:
        if not discussion:
            continue

        title = normalize(discussion.get("title"))
        body = normalize(discussion.get("body"))

        score, signals = score_text(title, body)

        usernames = set()

        author = discussion.get("author")
        if author and author.get("login"):
            usernames.add(author["login"])

        comments = (
            discussion.get("comments", {})
            .get("nodes", [])
        )

        for comment in comments:
            comment_author = comment.get("author")

            if comment_author and comment_author.get("login"):
                usernames.add(comment_author["login"])

            comment_body = normalize(comment.get("body"))

            extra_score, extra_signals = score_text(
                title,
                f"{body} {comment_body}",
            )

            score = max(score, extra_score)

            for signal in extra_signals:
                if signal not in signals:
                    signals.append(signal)

        if score <= 0:
            continue

        results.append(
            {
                "source": "discussion",
                "repository": f"{owner}/{name}",
                "number": discussion.get("number"),
                "title": title,
                "url": discussion.get("url"),
                "updated_at": discussion.get("updatedAt"),
                "score": score,
                "signals": sorted(signals),
                "public_users": sorted(usernames),
                "draft_reply": build_draft(
                    f"{owner}/{name}",
                    discussion.get("number"),
                    title,
                    discussion.get("url"),
                    signals,
                ),
            }
        )

    return sorted(
        results,
        key=lambda item: item["score"],
        reverse=True,
    )[:MAX_DISCUSSIONS_PER_REPOSITORY]


def collect_issues(
    github: Github,
    repository_name: str,
) -> list[dict[str, Any]]:
    repository = github.get_repo(repository_name)

    results = []

    issues = repository.get_issues(
        state="open",
        sort="updated",
        direction="desc",
    )

    for issue in issues:
        if issue.pull_request is not None:
            continue

        title = normalize(issue.title)
        body = normalize(issue.body)

        score, signals = score_text(title, body)

        if score <= 0:
            continue

        usernames = set()

        if issue.user and issue.user.login:
            usernames.add(issue.user.login)

        try:
            comments = issue.get_comments()

            for comment in comments:
                if comment.user and comment.user.login:
                    usernames.add(comment.user.login)

                comment_body = normalize(comment.body)

                extra_score, extra_signals = score_text(
                    title,
                    f"{body} {comment_body}",
                )

                score = max(score, extra_score)

                for signal in extra_signals:
                    if signal not in signals:
                        signals.append(signal)

        except GithubException:
            pass

        results.append(
            {
                "source": "issue",
                "repository": repository_name,
                "number": issue.number,
                "title": title,
                "url": issue.html_url,
                "updated_at": issue.updated_at.isoformat()
                if issue.updated_at
                else None,
                "score": score,
                "signals": sorted(signals),
                "public_users": sorted(usernames),
                "draft_reply": build_draft(
                    repository_name,
                    issue.number,
                    title,
                    issue.html_url,
                    signals,
                ),
            }
        )

        if len(results) >= MAX_ISSUES_PER_REPOSITORY:
            break

        time.sleep(0.15)

    return sorted(
        results,
        key=lambda item: item["score"],
        reverse=True,
    )


def collect_wiki_signal(repository_name: str) -> dict[str, Any] | None:
    """
    Wiki discovery is intentionally lightweight.

    GitHub's standard REST API does not expose wiki page content in the same
    way as Issues. We therefore record the public wiki endpoint as a
    discovery surface instead of cloning arbitrary repositories.
    """

    owner, name = repository_name.split("/", 1)

    wiki_url = f"https://github.com/{owner}/{name}/wiki"

    return {
        "source": "wiki",
        "repository": repository_name,
        "url": wiki_url,
        "note": (
            "Wiki endpoint discovered. Review pages manually for "
            "network/reliability signals."
        ),
    }


def main() -> None:
    token = os.getenv("GITHUB_TOKEN")

    if not token:
        raise RuntimeError(
            "GITHUB_TOKEN is required."
        )

    github = Github(token)

    OUTPUT_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    all_results: list[dict[str, Any]] = []
    wiki_results: list[dict[str, Any]] = []

    print("=" * 72)
    print("FAULTPLANE TRACTION HUNTER")
    print("=" * 72)
    print("Mode: discovery + scoring + draft generation")
    print("Automatic posting: DISABLED")
    print()

    for repository_name in TARGET_REPOSITORIES:
        print(f"[SCAN] {repository_name}")

        try:
            issue_results = collect_issues(
                github,
                repository_name,
            )

            all_results.extend(issue_results)

            for result in issue_results:
                print(
                    f"  ISSUE #{result['number']} "
                    f"score={result['score']} "
                    f"{result['title']}"
                )

            owner, name = repository_name.split("/", 1)

            try:
                discussion_results = collect_discussions(
                    token,
                    owner,
                    name,
                )

                all_results.extend(discussion_results)

                for result in discussion_results:
                    print(
                        f"  DISCUSSION #{result['number']} "
                        f"score={result['score']} "
                        f"{result['title']}"
                    )

            except Exception as exc:
                print(
                    f"  [WARN] Discussions unavailable: {exc}"
                )

            wiki = collect_wiki_signal(repository_name)

            if wiki:
                wiki_results.append(wiki)

        except GithubException as exc:
            print(
                f"  [WARN] GitHub API error: "
                f"{exc.status} {exc.data}"
            )

        except Exception as exc:
            print(f"  [WARN] Scan failed: {exc}")

        print()

    all_results.sort(
        key=lambda item: item.get("score", 0),
        reverse=True,
    )

    all_results = all_results[:MAX_TOTAL_RESULTS]

    users: dict[str, dict[str, Any]] = defaultdict(
        lambda: {
            "issues": [],
            "discussions": [],
            "repositories": set(),
        }
    )

    for result in all_results:
        for username in result.get("public_users", []):
            users[username]["repositories"].add(
                result["repository"]
            )

            if result["source"] == "issue":
                users[username]["issues"].append(
                    result["url"]
                )
            elif result["source"] == "discussion":
                users[username]["discussions"].append(
                    result["url"]
                )

    serializable_users = {}

    for username, data in users.items():
        serializable_users[username] = {
            "repositories": sorted(data["repositories"]),
            "issues": data["issues"],
            "discussions": data["discussions"],
        }

    payload = {
        "generated_at": time.strftime(
            "%Y-%m-%dT%H:%M:%SZ",
            time.gmtime(),
        ),
        "target_repositories": TARGET_REPOSITORIES,
        "total_matches": len(all_results),
        "matches": all_results,
        "public_users": serializable_users,
        "wiki_surfaces": wiki_results,
        "automatic_outreach": False,
    }

    JSON_OUTPUT.write_text(
        json.dumps(
            payload,
            indent=2,
            ensure_ascii=False,
        ),
        encoding="utf-8",
    )

    markdown = [
        "# FaultPlane Traction Hunter Report",
        "",
        f"Generated: {payload['generated_at']}",
        "",
        f"Matches: **{len(all_results)}**",
        "",
        "> Discovery only. No automatic comments, DMs, or mentions are sent.",
        "",
        "## Top Production Signals",
        "",
    ]

    for index, result in enumerate(all_results, 1):
        markdown.extend(
            [
                f"### {index}. {result['repository']} "
                f"#{result['number']}",
                "",
                f"**Score:** {result['score']}",
                "",
                f"**Title:** {result['title']}",
                "",
                f"**URL:** {result['url']}",
                "",
                "**Signals:** "
                + ", ".join(result["signals"]),
                "",
                "**Public users:** "
                + ", ".join(result["public_users"]),
                "",
                "#### Draft",
                "",
                "```text",
                result["draft_reply"],
                "```",
                "",
            ]
        )

    markdown.extend(
        [
            "## Wiki Surfaces",
            "",
        ]
    )

    for wiki in wiki_results:
        markdown.append(
            f"- {wiki['repository']}: {wiki['url']}"
        )

    MARKDOWN_OUTPUT.write_text(
        "\n".join(markdown),
        encoding="utf-8",
    )

    print("=" * 72)
    print(f"TOTAL MATCHES: {len(all_results)}")
    print(f"UNIQUE PUBLIC USERS: {len(serializable_users)}")
    print(f"JSON: {JSON_OUTPUT}")
    print(f"REPORT: {MARKDOWN_OUTPUT}")
    print("=" * 72)
    print()
    print("AUTOMATIC OUTREACH: DISABLED")
    print("HUMAN REVIEW REQUIRED")


if __name__ == "__main__":
    main()
