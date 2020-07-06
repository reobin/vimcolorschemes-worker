import datetime
import dateutil.parser as dparser
import os
import time

import db_service
import utils
import github
import printer
import request

MAX_IMAGE_COUNT = os.getenv("MAX_IMAGE_COUNT")
MAX_IMAGE_COUNT = int(MAX_IMAGE_COUNT) if MAX_IMAGE_COUNT is not None else 5

BUILD_WEBHOOK = os.getenv("BUILD_WEBHOOK")


def handler(event, context):
    start = time.time()

    printer.break_line()
    printer.break_line()

    last_import_at = db_service.get_last_import_at()

    printer.break_line()
    printer.break_line()

    repositories = github.search_repositories()

    for repository in repositories:
        owner_name = repository["owner"]["name"]
        name = repository["name"]

        is_new_repository, old_repository = db_service.is_repository_new(
            owner_name, name
        )

        repository["last_commit_at"] = github.get_last_commit_at(repository)
        last_commit_at = (
            dparser.parse(repository["last_commit_at"], fuzzy=True)
            if repository["last_commit_at"] is not None
            else None
        )

        fetch_images = (
            is_new_repository
            or last_import_at is None
            or last_commit_at > last_import_at
        )

        printer.info(
            "Images will be fetched" if fetch_images else "Image fetching skipped"
        )

        if fetch_images:
            image_urls = (
                list(
                    filter(
                        lambda url: request.is_image_url_valid(url),
                        old_repository["image_urls"],
                    )
                )
                if old_repository is not None
                else []
            )

            image_count_to_find = MAX_IMAGE_COUNT - len(image_urls)

            readme_image_urls = (
                utils.find_image_urls(
                    github.get_readme_file(repository), image_count_to_find, image_urls
                )
                if image_count_to_find > 0
                else []
            )

            image_urls = image_urls + readme_image_urls
            image_count_to_find = MAX_IMAGE_COUNT - len(image_urls)

            repository_image_urls = (
                github.list_repository_image_urls(
                    repository, image_count_to_find, image_urls
                )
                if image_count_to_find > 0
                else []
            )

            image_urls = image_urls + repository_image_urls

            repository["image_urls"] = list(map(utils.urlify, image_urls))

            printer.info(f"{len(readme_image_urls)} images found in readme")
            printer.info(f"{len(repository_image_urls)} images found in files")

        db_service.upsert_repository(repository)

        printer.break_line()
        printer.break_line()

    end = time.time()

    elapsed_time = end - start

    import_data = {
        "elapsed_time": elapsed_time,
        "import_at": datetime.datetime.now(datetime.timezone.utc),
    }

    db_service.create_import(import_data)

    if BUILD_WEBHOOK is not None:
        printer.break_line()
        printer.info("Starting website build")
        response = request.post(
            ***REMOVED***,
            is_json=False,
        )
        printer.info(f"Response: {response}")
        printer.break_line()

    return {"statusCode": 200, "body": f"Elapsed time: {elapsed_time}"}


if __name__ == "__main__":
    handler(None, None)
