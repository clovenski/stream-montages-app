using MontageJobExecutor.Models;

namespace MontageJobExecutor.Contracts
{
    public interface IMontageRepository
    {
        Task<CreateMontageResponse> CreateMontageAsync(CreateMontageRequest request);
    }

    public class CreateMontageRequest
    {
        public Montage Entity { get; set; }
    }
    
    public class CreateMontageResponse
    {
        public Montage CreatedEntity { get; set; }
    }
}
