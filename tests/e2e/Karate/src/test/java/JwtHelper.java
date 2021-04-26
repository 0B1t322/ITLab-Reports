import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTCreator;
import com.auth0.jwt.algorithms.Algorithm;
import jdk.nashorn.api.scripting.ScriptObjectMirror;

import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.Date;
import java.util.List;
import java.util.Map;

public class JwtHelper {
    public static String getAuthToken(
            String audience,
            String subject,
            Map<String, ScriptObjectMirror> claims) {
        try {
            Date now = new Date();
            Algorithm algorithm = Algorithm.HMAC256("test");
            JWTCreator.Builder builder = JWT.create()
                    .withIssuer("karate-tests")
                    .withAudience(audience)
                    .withSubject(subject)
                    .withExpiresAt(Date.from(LocalDateTime.now().plusDays(1).toInstant(ZoneOffset.UTC)));
            for (Map.Entry<String, ScriptObjectMirror> claim : claims.entrySet()) {
                builder.withArrayClaim(claim.getKey(), claim.getValue().to(String[].class));
            }
            return builder.sign(algorithm);
        } catch (Exception ex) {
            return ex.getMessage();
        }
    }
}
