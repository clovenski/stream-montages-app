using MontageJobExecutor.Contracts;

namespace MontageJobExecutor.Services
{
    public class MontageBuilder : IMontageBuilder
    {
        readonly IYouTubeAVStreamProvider _youTubeAVStreamProvider;
        readonly IClipService _clipService;
        readonly IVideoTranscriber _videoTranscriber;

        public MontageBuilder(
            IYouTubeAVStreamProvider youTubeAVStreamProvider,
            IClipService clipService,
            IVideoTranscriber videoTranscriber
        )
        {
            _youTubeAVStreamProvider = youTubeAVStreamProvider;
            _clipService             = clipService;
            _videoTranscriber        = videoTranscriber;
        }

        public async Task<BuildMontageResponse> Build(BuildMontageRequest request)
        {
            if (!request.Highlights.Any())
                throw new ArgumentException("Highlights arg must be non-empty");

            var clipPaths = new List<string>();
            foreach (var (highlight, idx) in request.Highlights.Select((it, i) => (it, i)))
            {
                // AV streams
                var streamUrls = await _youTubeAVStreamProvider.GetAVStreamUrlsAsync(new()
                {
                    VideoUrl = highlight.VideoURL,
                });
                var videoStreamUrl = streamUrls.VideoStreamUrl;
                var audioStreamUrl = streamUrls.AudioStreamUrl;

                // generating clips
                var clipPath = Path.Combine(request.OutputPathBase, $"{idx + 1:00}-montage-clip.mp4");
                clipPaths.Add(clipPath);
                var clipResult = await _clipService.GenerateClipAsync(new()
                {
                    VideoStreamUrl  = videoStreamUrl,
                    AudioStreamUrl  = audioStreamUrl,
                    DurationSeconds = (uint)highlight.DurationSeconds,
                    Timestamp       = $"{highlight.MinuteTimestamp}:00",
                    OutputPath      = clipPath,
                });
                Console.WriteLine($"Took {clipResult.Duration} to generate clip to path {clipResult.OutputPath}");
            }

            var montagePath = Path.Combine(request.OutputPathBase, "montage.mp4");
            if (clipPaths.Count > 1)
            {
                // join the clips together
                var concatResult = await _clipService.ConcatenateClipsAsync(new()
                {
                    Paths      = clipPaths,
                    OutputPath = montagePath,
                });
                Console.WriteLine($"Successfully created montage at {montagePath}");
            }
            else // 1 clip == montage
            {
                File.Copy(clipPaths.First(), montagePath);
            }

            // subtitles
            var transcription = await _videoTranscriber.TranscribeAsync(new()
            {
                SourcePath = montagePath,
            });
            Console.WriteLine($"Successfully transcribed video {montagePath} . Transcription written to {transcription.OutputPath}");

            // ??? somehow create new video file with subtitles overlaying it
            //     see https://www.bannerbear.com/blog/how-to-add-subtitles-to-a-video-file-using-ffmpeg/ for a start
            //     for now output montage and subtitles are fine for manual review, etc.
            
            return new()
            {
                VideoPath         = montagePath,
                TranscriptionPath = transcription.OutputPath,
            };
        }
    }
}
