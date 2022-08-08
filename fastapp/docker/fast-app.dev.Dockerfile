FROM python:3.9

ENV PORT=8001
EXPOSE ${PORT}
RUN mkdir /hyperpackage
WORKDIR /app
ADD docker/entrypoint-dev.sh ./entrypoint-dev.sh
RUN wget --quiet https://repo.anaconda.com/miniconda/Miniconda3-py39_4.10.3-Linux-x86_64.sh -O ./miniconda.sh && \
    bash ./miniconda.sh -b -p /opt/conda 
ENV PATH=/opt/conda/bin:$PATH
ADD environment.yml .
RUN conda env create -f environment.yml
ADD fastapp fastapp
ENTRYPOINT [ "./entrypoint-dev.sh" ]
