using MontageJobExecutor.Contracts;
using Xabe.FFmpeg;

namespace MontageJobExecutor.Services
{
    public class FFmpegClipService : IClipService
    {
        public async Task<GenerateClipResponse> GenerateClipAsync(GenerateClipRequest request)
        {
            var result = await FFmpeg.Conversions.New()
               .Start($"-ss {request.Timestamp} -i \"{request.VideoStreamUrl}\" -ss {request.Timestamp} -i \"{request.AudioStreamUrl}\" -map 0:v -map 1:a -t {request.DurationSeconds} -c:v libx264 -c:a aac {request.OutputPath} -y");

            return new()
            {
                OutputPath = request.OutputPath,
                Duration   = result.Duration,
            };
        }

        public async Task<ConcatenateClipsResponse> ConcatenateClipsAsync(ConcatenateClipsRequest request)
        {
            var concatListPath    = Path.Combine(Path.GetTempPath(), $"{Guid.NewGuid()}.txt");
            var concatListContent = string.Join(Environment.NewLine, request.Paths.Select(path => $"file '{Path.GetFullPath(path)}'"));
            File.WriteAllText(concatListPath, concatListContent);

            var result = await FFmpeg.Conversions.New()
                .Start($"-safe 0 -f concat -i {concatListPath} -c copy {request.OutputPath} -y");

            return new()
            {
                OutputPath = request.OutputPath,
            };
        }

        public async Task<GenerateAudioFileFromVideoResponse> GenerateAudioFileFromVideo(GenerateAudioFileFromVideoRequest request)
        {
            var outputPath = Path.Combine(
                Path.GetDirectoryName(request.SourcePath) ?? "",
                $"{Path.GetFileNameWithoutExtension(request.SourcePath)}.mp3"
            );

            var result = await FFmpeg.Conversions.New()
                .Start($"-i {request.SourcePath} {outputPath} -y");

            return new()
            {
                OutputPath = outputPath,
            };
        }
    }
}
