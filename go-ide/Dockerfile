FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install software-properties-common -y
RUN add-apt-repository ppa:neovim-ppa/stable \
 && apt-get update \
 && apt-get install neovim -y


RUN update-alternatives --install /usr/bin/vi vi /usr/bin/nvim 60 \
  && update-alternatives --config vi \
  && update-alternatives --install /usr/bin/vim vim /usr/bin/nvim 60 \
  && update-alternatives --config vim \
  && update-alternatives --install /usr/bin/editor editor /usr/bin/nvim 60 \
  && update-alternatives --config editor

RUN apt-get install curl git -y
RUN curl -o /tmp/go.tar.gz https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz \
 && tar -C /usr/local -xzvf /tmp/go.tar.gz \
 && rm -rf /tmp/go.tar.gz


RUN mkdir go
ENV GOPATH="/go:/app"
ENV GOBIN="/go/bin"
ENV PATH="$PATH:$GOPATH/bin:/usr/local/go/bin"

# vim-plug
RUN go get golang.org/x/tools/cmd/goimports
RUN go get -u github.com/nsf/gocode
RUN curl -fLo ~/.config/nvim/autoload/plug.vim --create-dirs \
  https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim


RUN echo "call plug#begin('~/.local/share/nvim/plugged')" \
    "\nPlug 'w0rp/ale'" \
    "\nPlug 'nsf/gocode', { 'rtp': 'nvim', 'do': '~/.config/nvim/plugged/gocode/nvim/symlink.sh' }" \
    "\ncall plug#end()" \
    "\n" \
    "\n" \
    "\nset expandtab" \
    "\nset ts=2" \
    "\nset sw=2" \
    "\ninoremap <C-Space> <C-x><C-o>" \
    "\ninoremap <C-@> <C-Space>" \
    "\n" \
    "\n" \
    > ~/.config/nvim/init.vim

RUN nvim +PlugInstall +qa

ADD GoLand-2017.3 /local/goland
RUN apt-get install x11-apps libxext-dev libxrender-dev libxslt1.1 libxtst-dev libgtk2.0-0 libcanberra-gtk-module -y

CMD ["/bin/bash"]
