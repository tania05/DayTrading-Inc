FROM ubuntu:16.04

# neo-vim
RUN apt-get update
RUN apt-get install software-properties-common -y
RUN add-apt-repository ppa:neovim-ppa/stable \
 && apt-get update \
 && apt-get install neovim -y
RUN apt-get install python-dev python-pip python3-dev python3-pip -y
RUN update-alternatives --install /usr/bin/vi vi /usr/bin/nvim 60 \
 && update-alternatives --config vi \
 && update-alternatives --install /usr/bin/vim vim /usr/bin/nvim 60 \
 && update-alternatives --config vim \
 && update-alternatives --install /usr/bin/editor editor /usr/bin/nvim 60 \
 && update-alternatives --config editor

# git
RUN add-apt-repository ppa:git-core/ppa \
 && apt-get update \
 && apt-get install git -y

# curl
RUN apt-get install curl -y

# stack
RUN echo "export PATH=~/.local/bin:\$PATH" >> ~/.bashrc
RUN curl -sSL https://get.haskellstack.org/ | sh

# vim-plug
RUN curl -fLo ~/.config/nvim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim

# build packages required for the vim configuration
RUN apt-get install libtinfo-dev -y
RUN stack build intero
RUN stack install hindent

RUN stack install hlint


RUN echo "call plug#begin('~/.local/share/nvim/plugged')" \
    "\nPlug 'parsonsmatt/intero-neovim'" \
    "\nPlug 'neovimhaskell/haskell-vim'" \
    "\nPlug 'alx741/vim-hindent'" \
    "\nPlug 'w0rp/ale'" \
    "\nPlug 'enomsg/vim-haskellConcealPlus'" \
    "\ncall plug#end()" \
    "\n" \
    "\n" \
    "\nlet g:ale_linters = {'haskell': ['hlint']}" \
    "\nlet g:intero_use_neomake = 0" \
    "\nlet g:hindent_on_save = 0" \
    "\n" \
    "\n" \
    "\nset expandtab" \
    "\nset ts=2" \
    "\nset sw=2" \
    "\n" \
    "\n" \
    > ~/.config/nvim/init.vim

# initalize autocomplete plugins for neo-vim
# this dosen't work sometimes....
RUN nvim +PlugInstall +qa

EXPOSE 8080
VOLUME ["/app"]
CMD ["/bin/bash"]

