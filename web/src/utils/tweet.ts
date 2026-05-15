// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025-2026 lin-snow

export type ParsedTweet = {
  url: string
  username: string
  statusId: string
}

const TWEET_HOST_PATTERN = /^(?:www\.)?(?:x|twitter)\.com$/i
const TWEET_PATH_PATTERN = /^\/([A-Za-z0-9_]{1,15})\/status(?:es)?\/(\d{1,32})\b/

export function parseTweetUrl(input: string): ParsedTweet | null {
  const raw = (input ?? '').trim()
  if (!raw) return null

  let parsed: URL
  try {
    parsed = new URL(raw)
  } catch {
    return null
  }

  if (!TWEET_HOST_PATTERN.test(parsed.hostname)) return null

  const match = parsed.pathname.match(TWEET_PATH_PATTERN)
  if (!match) return null

  const username = match[1]
  const statusId = match[2]
  return {
    url: `https://x.com/${username}/status/${statusId}`,
    username,
    statusId,
  }
}
