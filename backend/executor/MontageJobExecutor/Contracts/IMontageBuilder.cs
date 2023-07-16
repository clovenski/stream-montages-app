using MontageJobExecutor.Models;

namespace MontageJobExecutor.Contracts
{
    public interface IMontageBuilder
    {
        Task<BuildMontageResponse> Build(BuildMontageRequest request);
    }

    public class BuildMontageRequest
    {
        public IEnumerable<HighlightInfo> Highlights { get; set; }

        public string OutputPathBase { get; set; }
    }

    public class BuildMontageResponse
    {
        public string VideoPath { get; set; }

        public string TranscriptionPath { get; set; }
    }
}
