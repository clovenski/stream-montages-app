FROM mcr.microsoft.com/dotnet/sdk:6.0-alpine AS build-image

WORKDIR /app

COPY . .

RUN dotnet restore MontageJobExecutor.csproj
RUN dotnet publish MontageJobExecutor.csproj \
    -c Release -o out -r linux-musl-x64 --sc

# install dependencies
FROM mcr.microsoft.com/dotnet/runtime:6.0-alpine AS install-deps-image

WORKDIR /app/deps

RUN apk update

RUN apk add wget tar xz

# ffmpeg
RUN mkdir ffmpeg && cd ffmpeg && \
    wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && \
    # TODO: checksum?
    tar -xf ffmpeg-release-amd64-static.tar.xz && \
    rm ffmpeg-release-amd64-static.tar.xz && \
    cp -a ffmpeg-*/. . && \
    rm -rf ffmpeg-*

FROM mcr.microsoft.com/dotnet/runtime:6.0-alpine

COPY --from=build-image /app/out /app/out

RUN apk update

COPY --from=install-deps-image /app/deps /usr/local/bin

RUN ln -s /usr/local/bin/ffmpeg/ffmpeg /usr/bin/ffmpeg && \
    ln -s /usr/local/bin/ffmpeg/ffprobe /usr/bin/ffprobe

# yt-dlp
RUN apk add --no-cache yt-dlp

ENTRYPOINT [ "dotnet", "/app/out/MontageJobExecutor.dll" ]
