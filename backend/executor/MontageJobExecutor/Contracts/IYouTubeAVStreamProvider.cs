namespace MontageJobExecutor.Contracts
{
    public interface IYouTubeAVStreamProvider
    {
        Task<GetAVStreamUrlsResponse> GetAVStreamUrlsAsync(GetAVStreamUrlsRequest request);
    }

    public class GetAVStreamUrlsRequest
    {
        public string VideoUrl { get; set; }
    }

    public class GetAVStreamUrlsResponse
    {
        public string VideoStreamUrl { get; set; }

        public string AudioStreamUrl { get; set; }
    }
}
