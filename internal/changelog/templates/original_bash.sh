#!/bin/bash

# List of repositories to process
REPOS=("$REPO/myapp-ui" "$REPO/myapp-service")
# ref: https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional
conventional_commit_types="fix|feat|build|ci|docs|perf|refactor|revert|style|test"
# Global variables to collect raw changelog entries across all repositories
all_features=""
all_bugfixes=""
all_others=""

# Function to prompt for tags for a specific repository
get_tags_for_repo() {
    local repo_name="$1"
    local -n from_tag_ref=$2  # nameref for from_tag
    local -n to_tag_ref=$3    # nameref for to_tag
    
    echo "Enter tags for $repo_name:"
    read -p "From tag: " from_tag_ref
    read -p "To tag: " to_tag_ref
    
    # Validate tags exist
    pushd "$repo_name" > /dev/null || exit 1
    if ! git rev-parse "$from_tag_ref" > /dev/null 2>&1; then
        echo "Error: Tag $from_tag_ref does not exist in $repo_name"
        popd > /dev/null
        return 1
    fi
    if ! git rev-parse "$to_tag_ref" > /dev/null 2>&1; then
        echo "Error: Tag $to_tag_ref does not exist in $repo_name"
        popd > /dev/null
        return 1
    fi
    popd > /dev/null
    return 0
}

# Function to generate changelog between two tags for a single repository
generate_changelog() {
  repo_path="$1"
  from_tag="$2"
  to_tag="$3"

  echo ""
  echo "Generating changelog for repository: $repo_path"

  # Navigate to the repository
  pushd "$repo_path" > /dev/null || exit 1

  # Get commit logs between the tags without the date
  commits=$(git log "$from_tag..$to_tag" --oneline --pretty=format:"%s")

  while IFS= read -r commit; do
    # Regular expression to check for valid conventional commit formats
    if [[ ! "$commit" =~ ^($conventional_commit_types)\([A-Za-z0-9-]*\):|^($conventional_commit_types): ]]; then
      echo "Skipping commit: $commit"
      continue
    fi

    # Strip the scope in the conventional commit
    commit_message=$(echo "$commit" | sed -E "s/^($conventional_commit_types)\([A-Za-z0-9-]*\):/\1:/")
    # Strip leading and trailing space
    commit_message=$(echo "$commit_message" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    echo "Adding   commit: $commit_message"

    # Add repository identifier to the commit message
    repo_identifier=$(basename "$repo_path")
    commit_message="[$repo_identifier] $commit_message"

    # Categorize based on commit type
    if [[ "$commit_message" =~ ^.*fix: ]]; then
      all_bugfixes+="$commit_message"$'\n'
    elif [[ "$commit_message" =~ ^.*feat: ]]; then
      all_features+="$commit_message"$'\n'
    elif [[ "$commit_message" =~ ^.*(build:|chore:|ci:|docs:|perf:|refactor:|revert:|style:|test:) ]]; then
      all_others+="$commit_message"$'\n'
    fi
  done <<< "$commits"
  popd > /dev/null
}

# Function to transform raw changelog entries into a formatted list without prefix and trailing '- \n'
format_changelog_section() {
  raw_entries="$1"
  formatted_section=""
  while IFS= read -r entry; do
    # Skip empty entries
    if [ -n "$entry" ]; then
      # Extract repository identifier and remove the prefix type
      repo_id=$(echo "$entry" | grep -o '\[.*\]')
      entry_without_repo=$(echo "$entry" | sed -E "s/\[.*\][[:space:]]*($conventional_commit_types):[[:space:]]*//")
      formatted_section+="- ${repo_id} ${entry_without_repo}"$'\n'
    fi
  done <<< "$raw_entries"
  # Remove the last newline
  echo -n "$formatted_section" | sed '$ s/\n$//'
}

# Main logic to handle multiple repositories and merge changelog entries for all categories
main() {
  # Declare variables for tags
  local ui_from_tag ui_to_tag
  local service_from_tag service_to_tag
  
  # Get tags for UI repository
  if ! get_tags_for_repo "${REPOS[0]}" ui_from_tag ui_to_tag; then
    exit 1
  fi
  
  # Get tags for Service repository
  if ! get_tags_for_repo "${REPOS[1]}" service_from_tag service_to_tag; then
    exit 1
  fi

  # Generate changelog for UI
  generate_changelog "${REPOS[0]}" "$ui_from_tag" "$ui_to_tag"
  
  # Generate changelog for Service
  generate_changelog "${REPOS[1]}" "$service_from_tag" "$service_to_tag"

  # Prepare the combined changelog content
  final_changelog="## Changelog ($(date +%Y-%m-%d))\n"
  final_changelog+="Jira: $ui_from_tag → $ui_to_tag\n"
  final_changelog+="Service: $service_from_tag → $service_to_tag\n\n"
  
  if [ -n "$all_features" ]; then
    formatted_features=$(format_changelog_section "$all_features")
    final_changelog+="### Features\n$formatted_features\n\n"
  fi
  if [ -n "$all_bugfixes" ]; then
    formatted_bugfixes=$(format_changelog_section "$all_bugfixes")
    final_changelog+="### Bug Fixes\n$formatted_bugfixes\n\n"
  fi
  if [ -n "$all_others" ]; then
    formatted_others=$(format_changelog_section "$all_others")
    final_changelog+="### Other Changes\n$formatted_others\n\n"
  fi

  if [ -f CHANGELOG.md ]; then
    # prepend
    echo -e "$final_changelog\n$(cat CHANGELOG.md)" > CHANGELOG.md
  else
    # create
    echo -e "$final_changelog" > CHANGELOG.md
  fi
}

# -------------------------------------
# create_changelog.sh starts here

main