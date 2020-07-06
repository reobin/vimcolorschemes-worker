import re

import request


def find_image_urls(file_content, image_count_to_find, current_image_urls):
    image_url_regex = r"\b(https?:\/\/\S+(?:png|jpe?g|webp))\b"
    standard_image_urls = re.findall(image_url_regex, file_content)

    github_camo_url_regex = (
        r"\bhttps?:\/\/camo.githubusercontent.com(\/[0-9a-zA-Z]*)+\b"
    )
    github_camo_urls = re.findall(github_camo_url_regex, file_content)

    image_urls = standard_image_urls + github_camo_urls
    unique_image_urls = list(
        filter(lambda url: url not in current_image_urls, image_urls)
    )

    valid_image_urls = []
    index = 0

    while len(valid_image_urls) < image_count_to_find and index < len(
        unique_image_urls
    ):
        image_url = unique_image_urls[index]
        if request.is_image_url_valid(image_url):
            valid_image_urls.append(image_url)
        index = index + 1

    return valid_image_urls


def urlify(in_string):
    return "%20".join(in_string.split())
