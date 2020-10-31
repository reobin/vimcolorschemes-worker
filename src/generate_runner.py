import json
import os
import pathlib

from runner import Runner
from utils import printer

VIM_UTILS_PATH = "/usr/src/app/vim_utils"


class GenerateRunner(Runner):
    def run(self):
        set_up_vim_colors()

        repositories = self.database.get_valid_repositories()

        for repository in repositories:
            owner_name = repository["owner"]["name"]
            name = repository["name"]

            color_scheme_names = (
                repository["vim_color_scheme_names"]
                if "vim_color_scheme_names" in repository
                else None
            )
            if color_scheme_names is None or len(color_scheme_names) == 0:
                color_scheme_names = [name]

            reset_vimrc()
            install_color_scheme(owner_name, name, color_scheme_names)
            colors_data = get_colors_data(owner_name, name)
            if colors_data is not None:
                repository["colors"] = colors_data
                self.database.upsert_repository(repository)

        return {"previews_generated": 0}


def set_up_vim_colors():
    printer.info("Set up vim color syntax")

    os.system("mkdir -p ~/.vim/pack/vim-polyglot/start/vim-polyglot/")

    os.system(
        "git clone https://github.com/sheerun/vim-polyglot ~/.vim/pack/vim-polyglot/start/vim-polyglot"
    )


def reset_vimrc():
    os.system(f"cat {VIM_UTILS_PATH}/set_termguicolors.vim > ~/.vimrc")

    printer.info("Install vim color scheme preview generator")

    # TODO Maybe publish a vim plugin that we can install here instead of having it in the code
    os.system(f"cat {VIM_UTILS_PATH}/vcspg.vim >> ~/.vimrc")

    # Disable prompts on error
    os.system(f"echo 'set shortmess=a' >> ~/.vimrc")


def install_color_scheme(owner_name, name, color_scheme_names):
    printer.info(f"Install {owner_name}/{name} color scheme")

    os.system(f"mkdir -p ~/.vim/pack/{name}/start/{name}")

    os.system(
        f"git clone --depth 1 \
              https://github.com/{owner_name}/{name}.git \
              ~/.vim/pack/{name}/start/{name}"
    )

    printer.info(f"Set color scheme as {color_scheme_names[0]}")

    os.system(f"echo 'silent! colorscheme {color_scheme_names[0]}' >> ~/.vimrc")


def get_colors_data(owner_name, name):
    printer.info(f"Get colors for {owner_name}/{name}")

    file_path = f"{VIM_UTILS_PATH}/tmp_{owner_name}_{name}.json"

    printer.info(f"Create file at {file_path}")

    try:
        os.system(f"touch {file_path}")

        os.system(
            f'vim -c ":call WriteColorValues(\\"{file_path}\\")" {VIM_UTILS_PATH}/code_sample.vim -c ":q"'
        )

        content = pathlib.Path(file_path).read_text()
        data = json.loads(content)
        return data
    except Exception as error:
        printer.error(error)
        return None

    os.system(f"rm {file_path}")
