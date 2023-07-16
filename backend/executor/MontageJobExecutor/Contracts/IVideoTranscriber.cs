namespace MontageJobExecutor.Contracts
{
    public interface IVideoTranscriber
    {
        Task<TranscribeResponse> TranscribeAsync(TranscribeRequest request);
    }

    public class TranscribeRequest
    {
        public string SourcePath { get; set; }
    }

    public class TranscribeResponse
    {
        public string OutputPath { get; set; }
    }
}
