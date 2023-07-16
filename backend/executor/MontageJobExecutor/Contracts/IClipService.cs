namespace MontageJobExecutor.Contracts
{
    public interface IClipService
    {
        Task<GenerateClipResponse> GenerateClipAsync(GenerateClipRequest request);

        Task<ConcatenateClipsResponse> ConcatenateClipsAsync(ConcatenateClipsRequest request);

        Task<GenerateAudioFileFromVideoResponse> GenerateAudioFileFromVideo(GenerateAudioFileFromVideoRequest request);
    }

    public class GenerateClipRequest
    {
        /// <summary>
        /// Timestamp down to the second.
        /// i.e. 1:05:45 -> hour 1, minute 5, second 45
        /// </summary>
        public string Timestamp { get; set; }

        public string VideoStreamUrl { get; set; }

        public string AudioStreamUrl { get; set; }

        public uint DurationSeconds { get; set; }

        public string OutputPath { get; set; }
    }

    public class GenerateClipResponse
    {
        public string OutputPath { get; set; }

        public TimeSpan Duration { get; set; }
    }

    public class ConcatenateClipsRequest
    {
        public IEnumerable<string> Paths { get; set; }

        public string OutputPath { get; set; }
    }

    public class ConcatenateClipsResponse
    {
        public string OutputPath { get; set; }
    }

    public class GenerateAudioFileFromVideoRequest
    {
        public string SourcePath { get; set; }
    }

    public class GenerateAudioFileFromVideoResponse
    {
        public string OutputPath { get; set; }
    }
}
