using NYoutubeDL;
using MontageJobExecutor.Contracts;

namespace MontageJobExecutor.Services
{
    public class YouTubeAVStreamProvider : IYouTubeAVStreamProvider
    {
        readonly string? _exePath;

        public YouTubeAVStreamProvider(string? exePath)
        {
            _exePath = exePath;
        }

        public async Task<GetAVStreamUrlsResponse> GetAVStreamUrlsAsync(GetAVStreamUrlsRequest request)
        {
            var youtubeDl = new YoutubeDLP();
            if (!string.IsNullOrEmpty(_exePath))
            {
                youtubeDl.YoutubeDlPath = _exePath;
            }

            youtubeDl.Options.VideoFormatOptions.FormatAdvanced = "bestvideo[height<=480]+bestaudio/best[height<=480]";
            youtubeDl.Options.VerbositySimulationOptions.GetUrl = true;
            youtubeDl.Options.VideoFormatOptions.YoutubeSkipDashManifest = true;
            youtubeDl.VideoUrl = request.VideoUrl;
            youtubeDl.Options.FilesystemOptions.NoCacheDir = true;

            string? videoStreamUrl = null;
            string? audioStreamUrl = null;

            youtubeDl.StandardOutputEvent += (sender, output) =>
            {
                Console.WriteLine($"YTAVSP [{request.VideoUrl}]: {output}");
                if (videoStreamUrl == null)
                {
                    videoStreamUrl = output;
                }
                else if (audioStreamUrl == null)
                {
                    audioStreamUrl = output;
                }
            };
            youtubeDl.StandardErrorEvent += (sender, errorOutput) =>
            {
                Console.WriteLine($"YTAVSP [{request.VideoUrl}]: {errorOutput}");
            };

            await youtubeDl.DownloadAsync();

            Console.WriteLine($"Video stream url: {videoStreamUrl}");
            Console.WriteLine($"Audio stream url: {audioStreamUrl}");
            if (string.IsNullOrEmpty(videoStreamUrl) || string.IsNullOrEmpty(audioStreamUrl))
            {
                throw new Exception("Video or audio stream url was still empty.");
            }

            return new()
            {
                VideoStreamUrl = videoStreamUrl,
                AudioStreamUrl = audioStreamUrl,
            };
        }
    }
}
